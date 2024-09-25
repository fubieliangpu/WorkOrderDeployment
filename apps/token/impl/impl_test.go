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
	req := token.NewIssueTokenRequest("admin3", "43211")
	tk, err := serviceImpl.IssueToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tk)
}

func TestRevolkToken(t *testing.T) {
	req := token.NewRevolkTokenRequest("crpnii8b6bt4mn20gdvg", "crpnii8b6bt4mn20ge00")
	ins, err := serviceImpl.RevolkToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
