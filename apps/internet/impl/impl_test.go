package impl_test

import (
	"context"

	_ "github.com/fubieliangpu/WorkOrderDeployment/apps"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/internet"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"github.com/fubieliangpu/WorkOrderDeployment/test"
)

var (
	serviceImpl internet.Service
	ctx         = context.Background()
)

func init() {
	test.DevelopmentSetup()
	serviceImpl = ioc.Controller.Get(internet.AppName).(internet.Service)
}
