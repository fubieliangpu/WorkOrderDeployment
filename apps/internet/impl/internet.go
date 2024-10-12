package impl

import (
	"context"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/internet"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/exception"
)

// 冲突检测
func (i *NetProdDeplImpl) ConflictCheck(ctx context.Context, in *internet.DeploymentNetworkProductRequest) (internet.ConfigConflictStatus, error) {
	//校验请求完整性和合法性
	if err := in.Validate(); err != nil {
		return internet.CONFLICT, exception.ErrValidateFailed(err.Error())
	}
	//ping测即将要配置的IP是否通，看是否有IP冲突
	//首先要查询到指定IDC、属于指定设备层级的设备并登录设备
	req := rcdevice.NewQueryDeviceListRequest()
	req.IDC = in.Idc
	req.DeviceLevel = &in.AccessDeviceLevel
	deviceset, err := i.ctldevice.QueryDeviceList(ctx, req)
	if err != nil {
		return internet.CONFLICT, exception.ErrServerInternal(err.Error())
	}
	if len(deviceset.Items) == 0 {
		return internet.CONFLICT, internet.ErrNoDeviceInIdc
	}
	//当请求部署的位置在核心层，需要在对应的vpn-Instance下检查路由条目及在其下ping测，作为第一阶段检验
	//根据不同的接入方式及准备接入的设备层进行冲突检测，需要再添加一个辅助函数判定设备接入层，不同接入层有不同的部署逻辑，所以检测机制不同
	switch in.ConnectMethod {
	case internet.VRRP:
		err := in.BasicCheck(deviceset.Items[0])
		if err != nil {
			return internet.CONFLICT, err
		}
		err = in.BasicCheck(deviceset.Items[1])
		if err != nil {
			return internet.CONFLICT, err
		}
		return internet.FIT, nil

	case internet.SINGLE:
		err := in.BasicCheck(deviceset.Items[0])
		if err != nil {
			return internet.CONFLICT, err
		}
		return internet.FIT, nil

	//暂不支持跨机房
	case internet.STATIC_MASTER_BACKUP:
		err := in.BasicCheck(deviceset.Items[0])
		if err != nil {
			return internet.CONFLICT, err
		}
		err = in.BasicCheck(deviceset.Items[1])
		if err != nil {
			return internet.CONFLICT, err
		}
		return internet.FIT, nil

	case internet.SHAREGATEWAY:
		err := in.BasicCheck(deviceset.Items[0])
		if err != nil {
			return internet.CONFLICT, err
		}
		return internet.FIT, nil

	case internet.STATIC_LOADBALANCE:
		err := in.BasicCheck(deviceset.Items[0])
		if err != nil {
			return internet.CONFLICT, err
		}
		err = in.BasicCheck(deviceset.Items[1])
		if err != nil {
			return internet.CONFLICT, err
		}
		return internet.FIT, nil

	default:
		return internet.CONFLICT, exception.ErrDeviceAccessMothed("The ConnectMethod is not supported")
	}

}

// 业务配置下发与业务配置回收
func (i *NetProdDeplImpl) ConfigDeployment(ctx context.Context, in *internet.DeploymentNetworkProductRequest) (*internet.NetProd, error) {
	//校验与查询设备
	if err := in.Validate(); err != nil {
		return nil, exception.ErrValidateFailed(err.Error())
	}
	req := rcdevice.NewQueryDeviceListRequest()
	req.IDC = in.Idc
	req.DeviceLevel = &in.AccessDeviceLevel
	deviceset, err := i.ctldevice.QueryDeviceList(ctx, req)
	if err != nil {
		return nil, exception.ErrServerInternal(err.Error())
	}
	if deviceset.Total == 0 {
		return nil, internet.ErrNoDeviceInIdc
	}
	//由于先前已经判断过冲突问题，此处其他判断如分配哪个端口，是否配置携带vpn-instance由用户脚本中体现，同时会存在临时变化，程序不做判断，
	//安全起见,配置前需要人工审核脚本,脚本名称固定为configdeploymentA.txt configdeploymentA.txt,对应主备设备，如果没有备设备则只有configdeploymentA.txt
	configreq := rcdevice.NewChangeDeviceConfigRequest(deviceset.Items[0].Name)
	configreq.UserFile = "user.yaml"
	configreq.DeploymentRecord = "configrecord.txt"
	res := internet.NewNetProd()
	res.DeploymentNetworkProductRequest = in
	if in.ConnectMethod == internet.STATIC_LOADBALANCE || in.ConnectMethod == internet.STATIC_MASTER_BACKUP || in.ConnectMethod == internet.VRRP {
		configreq.DeviceConfigFile = "configdeploymentA.txt"
		if _, err := i.ctldevice.ChangeDeviceConfig(ctx, configreq); err != nil {
			res.Status = internet.FAIL
			return res, err
		}
		configreq.DeviceConfigFile = "configdeploymentB.txt"
		configreq.DeviceName = deviceset.Items[1].Name
		if _, err := i.ctldevice.ChangeDeviceConfig(ctx, configreq); err != nil {
			res.Status = internet.FAIL
			return res, err
		}
	} else if in.ConnectMethod == internet.SHAREGATEWAY || in.ConnectMethod == internet.SINGLE {
		configreq.DeviceConfigFile = "configdeploymentA.txt"
		if _, err := i.ctldevice.ChangeDeviceConfig(ctx, configreq); err != nil {
			res.Status = internet.FAIL
			return res, err
		}
	} else {
		res.Status = internet.FAIL
		return res, exception.ErrDeviceAccessMothed("The ConnectMethod is not supported")
	}
	res.Status = internet.SUCCESS
	return res, nil
}
