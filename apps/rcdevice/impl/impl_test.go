package impl_test

import (
	"context"

	_ "github.com/fubieliangpu/WorkOrderDeployment/apps"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/test"

	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
)

var (
	serviceImpl rcdevice.Service
	ctx         = context.Background()
)

func init() {
	test.DevelopmentSetup()
	serviceImpl = ioc.Controller.Get(rcdevice.AppName).(rcdevice.Service)
}
