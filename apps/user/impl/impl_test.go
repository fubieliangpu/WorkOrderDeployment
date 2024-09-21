package impl_test

import (
	"context"
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"github.com/fubieliangpu/WorkOrderDeployment/test"
)

var (
	serviceImpl user.Service
	ctx         = context.Background()
)

func init() {
	test.DevelopmentSetup()

	serviceImpl = ioc.Controller.Get(user.AppName).(user.Service)
}

func TestCreateVisitorUser(t *testing.T) {
	req := user.NewCreateUserRuquest()
	req.Username = "visitor1"
	req.Password = "123456"
	req.Role = user.ROLE_VISITOR
	ins, err := serviceImpl.CreateUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
