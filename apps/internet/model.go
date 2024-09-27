package internet

import (
	"encoding/json"

	"github.com/fubieliangpu/WorkOrderDeployment/common"
)

// 公网网络产品
type NetProd struct {
	*DeploymentNetworkProductRequest
	Status DeploymentResult `json:"status"`
}

func NewNetProd() *NetProd {
	return &NetProd{
		NewDeploymentNetworkProductRequest(),
		NewDeploymentResult(),
	}
}

func (req *NetProd) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}

type DeploymentNetworkProductRequest struct {
	//与客户对接的方式
	ConnectMethod ConMethod `json:"con_meth" validate:"required"`
	//客户与哪一层设备对接
	AccessDeviceLevel DeviceLevel `json:"dev_level" validate:"required"`
	//客户在哪个机房接入
	Idc string `json:"idc" validate:"required"`
	//客户接入的网络产品是哪个运营商的
	Operators string `json:"operators" validate:"required"`
	//为客户分配的业务IP
	IpAddr string `json:"ip_addr" validate:"required"`
	//与客户的互联IP段
	NeighborIp string `json:"neighbor_ip" validate:"required"`
}

func NewDeploymentNetworkProductRequest() *DeploymentNetworkProductRequest {
	return &DeploymentNetworkProductRequest{}
}

func (req *DeploymentNetworkProductRequest) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}

// 验证请求过来的字段是否有效
func (req *DeploymentNetworkProductRequest) Validate() error {
	return common.Validate(req)
}

type UndoDeviceConfigRequest struct {
	DeploymentNetworkProductRequest
}

func NewUndoDeviceConfigRequest() *UndoDeviceConfigRequest {
	return &UndoDeviceConfigRequest{}
}

func (req *UndoDeviceConfigRequest) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}

// 验证请求过来的字段是否有效
func (req *UndoDeviceConfigRequest) Validate() error {
	return common.Validate(req)
}
