package api

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/internet"
	"github.com/fubieliangpu/WorkOrderDeployment/conf"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
)

func init() {
	ioc.Api.Registry(internet.AppName, &InternetApiHandler{})
}

type InternetApiHandler struct {
	svc internet.Service
}

func (h *InternetApiHandler) Init() error {
	h.svc = ioc.Controller.Get(internet.AppName).(internet.Service)
	//注册Root Router
	subRouter := conf.C().Application.GinRootRouter().Group("internet")
	//Registory
	h.Registry(subRouter)
	return nil
}
