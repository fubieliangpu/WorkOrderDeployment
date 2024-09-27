package internet

type ConMethod int

const (
	VRRP ConMethod = iota + 1
	STATIC_MASTER_BACKUP
	SINGLE
	SHAREGATEWAY
	STATIC_LOADBALANCE
)

type DeploymentResult int

const (
	SUCCESS DeploymentResult = iota + 1
	FAIL
)

func NewDeploymentResult() DeploymentResult {
	return FAIL
}

type ConfigConflictStatus int

const (
	CONFLICT ConfigConflictStatus = iota + 1
	FIT
)
