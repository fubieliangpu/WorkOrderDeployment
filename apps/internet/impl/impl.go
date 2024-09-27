package impl

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/internet"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
)

func init() {
	ioc.Controller.Registry(internet.AppName, &NetProdDeplImpl{})
}

type NetProdDeplImpl struct {
	ctldevice rcdevice.Service
}

func (i *NetProdDeplImpl) Init() error {
	i.ctldevice = ioc.Controller.Get(rcdevice.AppName).(rcdevice.Service)
	return nil
}
