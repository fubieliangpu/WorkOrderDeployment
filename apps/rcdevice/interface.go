package rcdevice

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fubieliangpu/WorkOrderDeployment/common"
	"golang.org/x/crypto/ssh"
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
	//设备配置变更
	ChangeDeviceConfig(context.Context, *ChangeDeviceConfigRequest) (*Device, error)
	//设备配置或状态查询
	QueryDeviceConfig(context.Context, *QueryDeviceConfigRequest) (string, error)
}

type QueryDeviceListRequest struct {
	*common.PageRequest
	KeyWords string `json:"keywords"`
	Status   *Status
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
	DeviceName   string `json:"device_name" validate:"required"`
	DeviceConfig string `json:"device_config" validate:"required"`
}

func (req *ChangeDeviceConfigRequest) Validate() error {
	return common.Validate(req)
}

type Configinfo struct {
	Username   string
	Password   string
	Ip         string
	Port       string
	Protocol   string
	Configfile string
}

func NewConfiginfo() *Configinfo {
	return &Configinfo{
		Port:       "22",
		Protocol:   "tcp",
		Configfile: "config.txt",
	}
}

func SshConfigTool(cfi *Configinfo) {
	config := &ssh.ClientConfig{
		User: cfi.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfi.Password),
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
	session.Stdout = os.Stdout
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

type QueryDeviceConfigRequest struct {
	DeviceName string `json:"device_name" validate:"required"`
	Command    string `json:"command"`
}

func (req *QueryDeviceConfigRequest) Validate() error {
	return common.Validate(req)
}

func NewQueryDeviceConfigRequest(dsname string) *QueryDeviceConfigRequest {
	return &QueryDeviceConfigRequest{
		DeviceName: dsname,
	}
}
