package rcdevice

import "github.com/fubieliangpu/WorkOrderDeployment/common"

//设备
type Device struct {
	*common.DevMeta
	*CreateDeviceRequest
	*ChangedDeviceStatusRequest
}

//为新登录设备创建设备清单条目
type CreateDeviceRequest struct {
	Name       string `json:"name"`
	ServerAddr string `json:"server_addr"`
}

//检验是否已存在设备清单
type ChangedDeviceStatusRequest struct {
	CreateAt int64 `json:"create_at" gorm:"column:create_at"`
	//设备状态信息，已创建、已修改、已删除
	Status *Status `json:"status" gorm:"column:status"`
}
