package impl

import (
	"context"

	"dario.cat/mergo"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/common"
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
	//ins.SetIDC(in.Idc).SetBrand(in.Brand)
	ins.SetStatus(rcdevice.STATUS_CREATED)
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
		query = query.Where("idc LIKE ?", "%"+in.KeyWords+"%")
	}
	if in.Status != nil {
		query = query.Where("status = ?", *in.Status)
	}

	//Count,查询总数统计
	err := query.Count(&set.Total).Error
	if err != nil {
		return nil, err
	}

	//查询
	err = query.Order("change_at DESC").Limit(in.PageSize).Offset(in.Offset()).Find(&set.Items).Error
	if err != nil {
		return nil, err
	}
	return set, nil
}

// 设备条目更新
func (i *DeviceServiceImpl) UpdateDevice(ctx context.Context, in *rcdevice.UpdateDeviceRequest) (*rcdevice.Device, error) {
	//首先查询需要更新的对象
	ins, err := i.DescribeDevice(ctx, rcdevice.NewDescribeDeviceRequest(in.DeviceName))
	if err != nil {
		return nil, err
	}
	switch in.UpdateMode {
	case common.UPDATE_MODE_PUT: //如果是枚举类型，如Brand，不在枚举类型内则正常修改，但是值为0则报错Key: 'CreateDeviceRequest.Brand' Error:Field validation for 'Brand' failed on the 'required' tag
		ins.CreateDeviceRequest = in.CreateDeviceRequest
	case common.UPDATE_MODE_PATCH: //如果是枚举类型，如Brand，不在枚举类型内则正常修改,但是值为0则跳过
		err := mergo.MergeWithOverwrite(ins.CreateDeviceRequest, in.CreateDeviceRequest)
		if err != nil {
			return nil, err
		}
	}
	//更新字段校验
	if err := ins.CreateDeviceRequest.Validate(); err != nil {
		return nil, exception.ErrValidateFailed(err.Error())
	}

	//执行更新
	ins.ChangedDeviceStatusRequest.SetStatus(rcdevice.STATUS_MODIFIED)
	err = i.db.WithContext(ctx).Table("devices").Save(ins).Error
	if err != nil {
		return nil, err
	}
	return ins, nil
}

// 设备条目删除
func (i *DeviceServiceImpl) DeleteDevice(ctx context.Context, in *rcdevice.DeleteDeviceRequest) (*rcdevice.Device, error) {
	//首先查询需要删除的设备
	ins, err := i.DescribeDevice(ctx, rcdevice.NewDescribeDeviceRequest(in.DeviceName))
	if err != nil {
		return nil, err
	}

	err = i.db.WithContext(ctx).Table("devices").Where("name = ?", in.DeviceName).Delete(&rcdevice.Device{}).Error
	if err != nil {
		return nil, err
	}
	return ins, nil
}

// 设备配置修改
func (i *DeviceServiceImpl) ChangeDeviceConfig(ctx context.Context, in *rcdevice.ChangeDeviceConfigRequest) (*rcdevice.Device, error) {
	//首先要验证合法性
	if err := in.Validate(); err != nil {
		return nil, exception.ErrValidateFailed(err.Error())
	}
	//根据设备名查询设备地址端口等信息
	ins, err := i.DescribeDevice(ctx, rcdevice.NewDescribeDeviceRequest(in.DeviceName))
	if err != nil {
		return nil, err
	}

	//将查询到的IP和端口信息用于SSH配置下发
	ncfi := rcdevice.NewConfigInfo()
	ncfi.Ip, ncfi.Port, ncfi.Configfile, ncfi.Recordfile = ins.ServerAddr, ins.Port, in.DeviceConfigFile, in.DeploymentRecord
	//如果加载设备登录用户名密码错误，则抛出自定义错误
	nufi := rcdevice.NewDeviceUserInfo()
	ncfi.UserInfo, err = rcdevice.LoadUsernmPasswdFromYaml(in.UserFile, nufi)
	if err != nil {
		return nil, err
	}
	//核心操作，SSH登录设备并修改配置
	rcdevice.SshConfigTool(ncfi)
	return ins, nil
}
