package rcdevice

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fubieliangpu/WorkOrderDeployment/common"
	"github.com/fubieliangpu/WorkOrderDeployment/exception"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

const (
	AppName = "rcdevice"
)

type Service interface {
	//设备列表查询
	QueryDeviceList(context.Context, *QueryDeviceListRequest) (*DeviceSet, error)
	//设备详情
	DescribeDevice(context.Context, *DescribeDeviceRequest) (*Device, error)
	//设备创建
	CreateDevice(context.Context, *CreateDeviceRequest) (*Device, error)
	//设备更新
	UpdateDevice(context.Context, *UpdateDeviceRequest) (*Device, error)
	//设备删除
	DeleteDevice(context.Context, *DeleteDeviceRequest) (*Device, error)
	//设备配置变更或查询或状态查询
	ChangeDeviceConfig(context.Context, *ChangeDeviceConfigRequest) (*Device, error)
	// //设备配置或状态查询
	// QueryDeviceConfig(context.Context, *QueryDeviceConfigRequest) (string, error)
}

type QueryDeviceListRequest struct {
	*common.PageRequest
	IDC         string              `json:"idc"`
	DeviceLevel *common.DeviceLevel `json:"device_level"`
	Status      *Status             `json:"status"`
}

func (c *QueryDeviceListRequest) SetStatus(v Status) {
	c.Status = &v
}

func NewQueryDeviceListRequest() *QueryDeviceListRequest {
	return &QueryDeviceListRequest{
		PageRequest: common.NewPageRequest(),
	}
}

type DescribeDeviceRequest struct {
	DeviceName string `json:"device_name"`
}

func NewDescribeDeviceRequest(dsname string) *DescribeDeviceRequest {
	return &DescribeDeviceRequest{
		DeviceName: dsname,
	}
}

type UpdateDeviceRequest struct {
	DeviceName           string `json:"device_name"`
	*CreateDeviceRequest `validate:"required"`
	UpdateMode           common.UPDATE_MODE `json:"update_mode"`
}

func (req *UpdateDeviceRequest) Validate() error {
	return common.Validate(req)
}

func NewUpdateDeviceRequest(dsname string) *UpdateDeviceRequest {
	return &UpdateDeviceRequest{
		DeviceName:          dsname,
		CreateDeviceRequest: NewCreateDeviceRequest(),
		UpdateMode:          common.UPDATE_MODE_PUT,
	}
}

type DeleteDeviceRequest struct {
	DeviceName string `json:"device_name"`
}

func NewDeleteDeviceRequest(dsname string) *DeleteDeviceRequest {
	return &DeleteDeviceRequest{
		DeviceName: dsname,
	}
}

type ChangeDeviceConfigRequest struct {
	DeviceName       string `json:"device_name" validate:"required"`
	DeviceConfigFile string `json:"device_config_file" validate:"required"`
	UserFile         string `json:"user_file" validate:"required"`
	DeploymentRecord string `json:"deployment_record" validate:"required"`
}

func (req *ChangeDeviceConfigRequest) Validate() error {
	return common.Validate(req)
}

type DeviceUserInfo struct {
	Username string `toml:"username" yaml:"username" json:"username"`
	Password string `toml:"password" yaml:"password" json:"password"`
}

func NewDeviceUserInfo() *DeviceUserInfo {
	return &DeviceUserInfo{}
}

type ConfigInfo struct {
	UserInfo   *DeviceUserInfo `toml:"userinfo" yaml:"userinfo" json:"userinfo"`
	Ip         string
	Port       string
	Protocol   string
	Configfile string
	Recordfile string
}

func NewConfigInfo() *ConfigInfo {
	return &ConfigInfo{
		Protocol: "tcp",
	}
}

// 从yaml文件读取登录设备的用户名密码
func LoadUsernmPasswdFromYaml(userpath string, deviceuserinfo *DeviceUserInfo) (*DeviceUserInfo, error) {
	content, err := os.ReadFile(userpath)
	if err != nil {
		return nil, exception.ErrOpenFileFailed(err.Error())
	}

	err = yaml.Unmarshal(content, deviceuserinfo)
	if err != nil {
		return nil, exception.ErrParseFileFailed(err.Error())
	}

	return deviceuserinfo, nil
}

func SshConfigTool(cfi *ConfigInfo) {
	config := &ssh.ClientConfig{
		User: cfi.UserInfo.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfi.UserInfo.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	svr := fmt.Sprintf("%s:%s", cfi.Ip, cfi.Port)
	client, err := ssh.Dial(cfi.Protocol, svr, config)
	if err != nil {
		log.Fatal("Fail to dial: ", err)
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Fail to dial: ", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("linux", 200, 200, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	fdd, _ := os.Open(cfi.Configfile)
	defer fdd.Close()
	//每一次操作都留记录，输出到文件中
	fresult, err := os.OpenFile(cfi.Recordfile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer fresult.Close()

	session.Stdout = fresult
	session.Stdin = fdd
	session.Stderr = os.Stderr
	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}

	err = session.Wait()
	if err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
}

func NewChangeDeviceConfigRequest(dsname string) *ChangeDeviceConfigRequest {
	return &ChangeDeviceConfigRequest{
		DeviceName: dsname,
	}
}

// 似乎不需要特别定义查询，复用配置修改就可以实现，后面完善sshtool的返回类型和输出就行
// type QueryDeviceConfigRequest struct {
// 	DeviceName  string `json:"device_name" validate:"required"`
// 	CommandFile string `json:"command"`
// }

// func (req *QueryDeviceConfigRequest) Validate() error {
// 	return common.Validate(req)
// }

// func NewQueryDeviceConfigRequest(dsname string) *QueryDeviceConfigRequest {
// 	return &QueryDeviceConfigRequest{
// 		DeviceName: dsname,
// 	}
// }
