package rcdevice

import (
	"context"

	"github.com/fubieliangpu/WorkOrderDeployment/common"
)

type Service interface {
	//设备列表查询
	QueryDeviceList(context.Context)
	//设备详情
	//设备创建
	//设备更新
	//设备删除
	//设备配置变更
	//设备配置或状态查询
}

type QueryDeviceListRequest struct {
	*common.PageRequest
	KeyWords string `json:"keywords"`
	Status   *Status
}
