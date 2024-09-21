package impl_test

import (
	"context"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/token"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"github.com/fubieliangpu/WorkOrderDeployment/test"
)

var (
	//申明被测试对象
	serviceImpl token.Service
	ctx         = context.Background()
)

func init() {
	test.DevelopmentSetup()
	serviceImpl = ioc.Controller.Get(token.AppName).(token.Service)
}
