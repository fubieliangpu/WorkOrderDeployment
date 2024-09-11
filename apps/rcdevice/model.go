package rcdevice

import (
	"encoding/json"
	"time"

	"github.com/fubieliangpu/WorkOrderDeployment/common"
)

// 设备
type Device struct {
	*common.DevMeta
	*CreateDeviceRequest
	*ChangedDeviceStatusRequest
}

func NewDevice(svridc string) *Device {
	return &Device{
		common.NewDevMeta(),
		NewCreateDeviceRequest(),
		NewChangedDeviceStatusRequest(),
	}
}

func (req *Device) SetIDC(svridc string) *Device {
	req.DevMeta.IDC = svridc
	return req
}

func (req *Device) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}

// 为新登录设备创建设备清单条目
type CreateDeviceRequest struct {
	Name       string `json:"name" gorm:"column:name" validate:"required"`
	ServerAddr string `json:"server_addr" gorm:"column:server_addr" validate:"required"`
	Port       string `json:"port" gorm:"column:port" validate:"required"`
}

func NewCreateDeviceRequest() *CreateDeviceRequest {
	return &CreateDeviceRequest{}
}

func (req *CreateDeviceRequest) SetDevice(svrname, svraddr, svrport string) *CreateDeviceRequest {
	req.Name, req.ServerAddr, req.Port = svrname, svraddr, svrport
	return req
}

// 验证请求数据
func (req *CreateDeviceRequest) Validate() error {
	return common.Validate(req)
}

func (req *CreateDeviceRequest) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}

// 检验是否已存在设备清单
type ChangedDeviceStatusRequest struct {
	ChangedAt int64 `json:"create_at" gorm:"column:create_at"`
	//设备状态信息，创建、修改、删除
	Status *Status `json:"status" gorm:"column:status"`
}

func NewChangedDeviceStatusRequest() *ChangedDeviceStatusRequest {
	return &ChangedDeviceStatusRequest{}
}

func (req *ChangedDeviceStatusRequest) SetStatus(s Status) *ChangedDeviceStatusRequest {
	req.Status = &s
	switch *req.Status {
	case STATUS_CREATED:
		req.ChangedAt = time.Now().Unix()
	case STATUS_MODIFIED:
		req.ChangedAt = time.Now().Unix()
	case STATUS_DELETED:
		req.ChangedAt = time.Now().Unix()
	}
	return req
}

func (req *ChangedDeviceStatusRequest) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}

// 设备清单列表
type DeviceSet struct {
	Total int64     `json:"total"`
	Items []*Device `json:"item"`
}

func NewDeviceSet() *DeviceSet {
	return &DeviceSet{
		Items: []*Device{},
	}
}

func (req *DeviceSet) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}
