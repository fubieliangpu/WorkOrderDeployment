package api

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
	"github.com/fubieliangpu/WorkOrderDeployment/conf"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
)

func init() {
	ioc.Api.Registry(user.AppName, &UserApiHandler{})
}

type UserApiHandler struct {
	//依赖User接口
	userapi user.Service
}

func (h *UserApiHandler) Init() error {
	h.userapi = ioc.Controller.Get(user.AppName).(user.Service)
	subRouter := conf.C().Application.GinRootRouter().Group("users")
	h.Registry(subRouter)
	return nil
}
