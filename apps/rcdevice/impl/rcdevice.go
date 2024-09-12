package impl

import (
	"context"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
)

// 设备创建
func (i *DeviceServiceImpl) CreateDevice(ctx context.Context, in *rcdevice.CreateDeviceRequest) (*rcdevice.Device, error) {
	//1.验证请求参数
	if err := in.Validate(); err != nil {
		return nil, err
	}
}
