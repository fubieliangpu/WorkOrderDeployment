package impl_test

import (
	"context"
	"testing"

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

func TestIssueToken(t *testing.T) {
	req := token.NewIssueTokenRequest("Admin1", "65432111")
	tk, err := serviceImpl.IssueToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tk)
}

func TestRevolkToken(t *testing.T) {
	req := token.NewRevolkTokenRequest("crntklp14uflk1303k20", "crntklp14uflk1303k2g")
	ins, err := serviceImpl.RevolkToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
