package dci

import (
	"encoding/json"

	"github.com/fubieliangpu/WorkOrderDeployment/common"
)

// DCI产品
type DCI struct {
	*CreateDCIRequest
}

func NewDCI() *DCI {
	return &DCI{
		NewCreateDCIRequest(),
	}
}

func (req *DCI) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}

// 创建DCI产品请求
type CreateDCIRequest struct {
	//在哪个机房部署
	Idc string `json:"idc" validate:"required"`
	//vni的号码，需要保证在此机房不重复
	Vni int `json:"vni" validate:"required"`
	//vlan号码
	Vlan int `json:"vlan" validate:"required"`
	//vsi的名字，需要在设备上不重复
	Vname string `json:"vname" validate:"required"`
	//聚合口端口号，需要在设备上不重复，如果不使用聚合口对接，则为0
	BridgePort string `json:"bridge_port" validate:"required"`
	//需求对接物理口数量，需要检测设备上空闲物理口的数量
	PortNumber int `json:"port_number" validate:"required"`
	//Tunnel目标的设备名，主要用来检查有没有已经存在的通道
	DestDevName string `json:"dest_dev_name" validate:"required"`
	//源设备名
	DevName string `json:"dev_name" validate:"required"`
}

func NewCreateDCIRequest() *CreateDCIRequest {
	return &CreateDCIRequest{}
}

// 验证客户请求的数据完整性
func (req *CreateDCIRequest) Validate() error {
	return common.Validate(req)
}

func (req *CreateDCIRequest) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}
