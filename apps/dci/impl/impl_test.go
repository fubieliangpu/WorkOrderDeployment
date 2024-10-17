package impl_test

import (
	"context"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/dci"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"github.com/fubieliangpu/WorkOrderDeployment/test"
)

var (
	serviceImpl dci.Service
	ctx         = context.Background()
)

func init() {
	test.DevelopmentSetup()
	serviceImpl = ioc.Controller.Get(dci.AppName).(dci.Service)
}
