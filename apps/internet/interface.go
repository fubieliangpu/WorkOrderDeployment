package internet

import (
	"context"
)

const (
	AppName = "internet"
)

type Service interface {
	//冲突检测
	VrrpConflictCheck(context.Context, *DeploymentVRRP) (ConfigConflictStatus, error)
	DoubleStaticConflictCheck(context.Context, *DeploymentDoubleStatic) (ConfigConflictStatus, error)
	SingleConflictCheck(context.Context, *DeploymentSingle) (ConfigConflictStatus, error)
	//配置下发与业务配置回收
	VrrpDeployment(context.Context, *DeploymentVRRP) (DeploymentResult, error)
	DoubleStaticDeployment(context.Context, *DeploymentDoubleStatic) (DeploymentResult, error)
	SingleDeployment(context.Context, *DeploymentSingle) (DeploymentResult, error)
}
