package impl_test

import (
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
)

func TestCreateDevice(t testing.T) {
	req := rcdevice.NewCreateDeviceRequest()
	req.Name = "RJJD-6890"
	req.ServerAddr = "192.168.101.220"
	req.Port = "22"
	req.Brand = 1
	req.Idc = "RJJD"
	ins, err := serviceImpl.CreateDevice(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
