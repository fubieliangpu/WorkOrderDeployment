package impl

import (
	"context"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/internet"
	"github.com/fubieliangpu/WorkOrderDeployment/exception"
)

// 冲突检测
func (i *NetProdDeplImpl) ConflictCheck(ctx context.Context, in *internet.DeploymentNetworkProductRequest) (internet.ConfigConflictStatus, error) {
	//校验请求完整性和合法性
	if err := in.Validate(); err != nil {
		return 1, exception.ErrValidateFailed(err.Error())
	}
	//校验ping测即将要配置的IP是否通，看是否有IP冲突
	//根据不同的接入方式及准备接入的设备层进行冲突检测，需要再添加一个辅助函数判定设备接入层，不同接入层有不同的部署逻辑，所以检测机制不同
	switch in.ConnectMethod {
	case internet.VRRP:
	case internet.SINGLE:
	case internet.STATIC_MASTER_BACKUP:
	case internet.SHAREGATEWAY:
	case internet.STATIC_LOADBALANCE:
	default:
	}

	return 1, nil
}

// 业务配置下发
func (i *NetProdDeplImpl) ConfigDeployment(ctx context.Context, in *internet.DeploymentNetworkProductRequest) (*internet.NetProd, error) {
	return nil, nil
}

// 业务配置回收
func (i *NetProdDeplImpl) ConfigRevoke(ctx context.Context, in *internet.UndoDeviceConfigRequest) (*internet.NetProd, error) {
	return nil, nil
}

//定义一个函数用于帮助冲突检测
