package impl_test

import (
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/internet"
	"github.com/fubieliangpu/WorkOrderDeployment/common"
)

func TestConflictCheck(t *testing.T) {
	req := internet.NewDeploymentNetworkProductRequest()
	req.AccessDeviceLevel = common.CORE
	req.ConfigRevoke = 1
	req.ConnectMethod = internet.STATIC_LOADBALANCE
	req.Idc = "HYC"
	req.IpAddr = "120.133.141.0"
	req.IpMask = "27"
	req.NeighborIp = "172.18.1.0"
	req.NeighborMask = "30"
	req.Operators = "LY-CT"
	ins, err := serviceImpl.ConflictCheck(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestConfigDeployment(t *testing.T) {
	req := internet.NewDeploymentNetworkProductRequest()
	req.AccessDeviceLevel = common.CORE
	req.ConfigRevoke = 1
	req.ConnectMethod = internet.STATIC_LOADBALANCE
	req.Idc = "HYC"
	req.IpAddr = "120.133.141.0"
	req.IpMask = "27"
	req.NeighborIp = "172.18.1.0"
	req.NeighborMask = "30"
	req.Operators = "LY-CT"
	ins, err := serviceImpl.ConfigDeployment(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
