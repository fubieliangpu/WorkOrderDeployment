package internet

import (
	"encoding/json"

	"github.com/fubieliangpu/WorkOrderDeployment/common"
)

// 业务和互联IP的掩码及运营商信息
type IpMaskOpt struct {
	//客户接入的网络产品是哪个运营商的
	Operators string `json:"operators" validate:"required"`
	//为客户分配的业务IP
	IpAddr string `json:"ip_addr" validate:"required"`
	//业务IP的掩码
	IpMask string `json:"ip_mask" validate:"required"`
	//与客户的互联IP段
	NeighborIp string `json:"neighbor_ip" validate:"required"`
	//与客户的互联IP段掩码
	NeighborMask string `json:"neighbor_maks" validate:"required"`
}

func NewIpMaskOpt() *IpMaskOpt {
	return &IpMaskOpt{}
}

func (req *IpMaskOpt) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}

func (req *IpMaskOpt) Validate() error {
	return common.Validate(req)
}

// 双设备VRRP
type DeploymentVRRP struct {
	//准备部署的VRID
	Vrid uint8 `json:"vrid" validate:"required"`
	//VRRP主设备名
	MasterDevName string `json:"master_dev" validate:"required"`
	//VRRP备设备名
	BackupDevName string     `json:"backup_dev" validate:"required"`
	Detail        *IpMaskOpt `json:"detail" validate:"required"`
}

func NewDeploymentVRRP() *DeploymentVRRP {
	return &DeploymentVRRP{
		Detail: NewIpMaskOpt(),
	}
}

func (req *DeploymentVRRP) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}

func (req *DeploymentVRRP) Validate() error {
	return common.Validate(req)
}

// 双设备静态部署
type DeploymentDoubleStatic struct {
	//第一台设备名
	FirstDevName string `json:"first_dev" validate:"required"`
	//第二台设备名
	SecondDevName string     `json:"second_dev" validate:"required"`
	Detail        *IpMaskOpt `json:"detail" validate:"required"`
}

func NewDeploymentDoubleStatic() *DeploymentDoubleStatic {
	return &DeploymentDoubleStatic{
		Detail: NewIpMaskOpt(),
	}
}

func (req *DeploymentDoubleStatic) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}

func (req *DeploymentDoubleStatic) Validate() error {
	return common.Validate(req)
}

// 单设备部署
type DeploymentSingle struct {
	//设备名
	DevName string     `json:"dev_name" validate:"required"`
	Detail  *IpMaskOpt `json:"detail" validate:"required"`
}

func NewDeploymentSingle() *DeploymentSingle {
	return &DeploymentSingle{
		Detail: NewIpMaskOpt(),
	}
}

func (req *DeploymentSingle) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}

func (req *DeploymentSingle) Validate() error {
	return common.Validate(req)
}
