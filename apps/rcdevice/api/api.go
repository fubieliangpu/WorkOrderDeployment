package api

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/conf"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
)

func init() {
	ioc.Api.Registry(rcdevice.AppName, &RcDeviceApiHandler{})
}

type RcDeviceApiHandler struct {
	svc rcdevice.Service
}

func (h *RcDeviceApiHandler) Init() error {
	h.svc = ioc.Controller.Get(rcdevice.AppName).(rcdevice.Service)
	//注册Root Router
	subRouter := conf.C().Application.GinRootRouter().Group("rcdevice")
	h.Registry(subRouter)
	return nil
}
