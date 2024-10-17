package api

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/dci"
	"github.com/fubieliangpu/WorkOrderDeployment/conf"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
)

func init() {
	ioc.Api.Registry(dci.AppName, &DCIApiHandler{})
}

type DCIApiHandler struct {
	svc dci.Service
}

func (h *DCIApiHandler) Init() error {
	h.svc = ioc.Controller.Get(dci.AppName).(dci.Service)
	//注册Root Router
	subRouter := conf.C().Application.GinRootRouter().Group("dci")
	//Registory
	h.Registry(subRouter)
	return nil
}
