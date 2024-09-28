package internet

import (
	"context"
)

const (
	AppName = "internet"
)

type Service interface {
	//冲突检测
	ConflictCheck(context.Context, *DeploymentNetworkProductRequest) (ConfigConflictStatus, error)
	//业务配置下发
	ConfigDeployment(context.Context, *DeploymentNetworkProductRequest) (*NetProd, error)
	//业务配置回收
	ConfigRevoke(context.Context, *UndoDeviceConfigRequest) (*NetProd, error)
}
