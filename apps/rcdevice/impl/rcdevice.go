package impl

import (
	"context"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/exception"
)

// 设备条目信息查询
func (i *DeviceServiceImpl) DescribeDevice(ctx context.Context, in *rcdevice.DescribeDeviceRequest) (*rcdevice.Device, error) {
	ins := rcdevice.NewDevice()
	err := i.db.WithContext(ctx).Where("name = ?", in.DeviceName).First(ins).Error
	if err != nil {
		return nil, err
	}
	return ins, nil
}

// 设备创建
func (i *DeviceServiceImpl) CreateDevice(ctx context.Context, in *rcdevice.CreateDeviceRequest) (*rcdevice.Device, error) {
	//1.验证请求参数
	if err := in.Validate(); err != nil {
		return nil, exception.ErrValidateFailed(err.Error())
	}
	//2.构造实例对象
	ins := rcdevice.NewDevice()
	ins.CreateDeviceRequest = in

	//3.入库返回
	err := i.db.WithContext(ctx).Create(ins).Error
	if err != nil {
		return nil, err
	}
	return ins, nil
}

// 设备列表查询
func (i *DeviceServiceImpl) QueryDeviceList(ctx context.Context, in *rcdevice.QueryDeviceListRequest) (*rcdevice.DeviceSet, error) {
	set := rcdevice.NewDeviceSet()
	//查找数据库，构造查询条件
	query := i.db.WithContext(ctx).Table("devices")
	if in.KeyWords != "" {
		query = query.Where("title LIKE ?", "%"+in.KeyWords+"%")
	}
	if in.Status != nil {
		query = query.Where("status = ?", *in.Status)
	}
}
