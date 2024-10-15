package impl

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/dci"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
)

func init() {
	ioc.Controller.Registry(dci.AppName, &DCIDeplImpl{})
}

type DCIDeplImpl struct {
	ctldevice rcdevice.Service
}

func (i *DCIDeplImpl) Init() error {
	i.ctldevice = ioc.Controller.Get(rcdevice.AppName).(rcdevice.Service)
	return nil
}
