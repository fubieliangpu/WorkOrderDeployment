package impl_test

import (
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/internet"
)

// func TestConflictCheck(t *testing.T) {
// 	req := internet.NewDeploymentNetworkProductRequest()
// 	req.AccessDeviceLevel = common.CORE
// 	req.ConfigRevoke = 1
// 	req.ConnectMethod = internet.STATIC_LOADBALANCE
// 	req.Idc = "HYC"
// 	req.IpAddr = "120.133.141.0"
// 	req.IpMask = "27"
// 	req.NeighborIp = "172.18.1.0"
// 	req.NeighborMask = "30"
// 	req.Operators = "LY-CT"
// 	ins, err := serviceImpl.ConflictCheck(ctx, req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(ins)
// }

// func TestConfigDeployment(t *testing.T) {
// 	req := internet.NewDeploymentNetworkProductRequest()
// 	req.AccessDeviceLevel = common.CORE
// 	req.ConfigRevoke = 1
// 	req.ConnectMethod = internet.STATIC_LOADBALANCE
// 	req.Idc = "HYC"
// 	req.IpAddr = "120.133.141.0"
// 	req.IpMask = "27"
// 	req.NeighborIp = "172.18.1.0"
// 	req.NeighborMask = "30"
// 	req.Operators = "LY-CT"
// 	ins, err := serviceImpl.ConfigDeployment(ctx, req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(ins)
// }

func TestVrrpConflictCheck(t *testing.T) {
	req := internet.NewDeploymentVRRP()
	req.MasterDevName = "SC-CDJF-IDC2-D13-X-S8861-01"
	req.BackupDevName = "SC-CDJF-IDC2-D13-X-S8861-02"
	req.Vrid = 3
	req.Detail.Operators = "CD-CT"
	req.Detail.IpAddr = "112.45.33.0"
	req.Detail.IpMask = "27"
	req.Detail.NeighborIp = "172.18.15.0"
	req.Detail.NeighborMask = "29"
	ins, err := serviceImpl.VrrpConflictCheck(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestDoubleStaticConflictCheck(t *testing.T) {
	req := internet.NewDeploymentDoubleStatic()
	req.FirstDevName = "FS-ZHXC-3L-HH08-04-X-S9850-03"
	req.SecondDevName = "FS-ZHXC-3L-HH08-05-X-S9850-04"
	req.Detail.Operators = "CD-CT"
	req.Detail.IpAddr = "112.45.34.0"
	req.Detail.IpMask = "27"
	req.Detail.NeighborIp = "172.18.15.0"
	req.Detail.NeighborMask = "29"
	ins, err := serviceImpl.DoubleStaticConflictCheck(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestSingleConflictCheck(t *testing.T) {
	req := internet.NewDeploymentSingle()
	req.DevName = "FS-ZHXC-3L-HH08-04-X-S9850-03"
	req.Detail.Operators = "CD-CMCC"
	req.Detail.IpAddr = "112.45.29.0"
	req.Detail.IpMask = "27"
	req.Detail.NeighborIp = "172.18.15.0"
	req.Detail.NeighborMask = "30"
	ins, err := serviceImpl.SingleConflictCheck(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestVrrpDeployment(t *testing.T) {
	req := internet.NewDeploymentVRRP()
	req.MasterDevName = "FS-ZHXC-3L-HH08-04-X-S9850-03"
	req.BackupDevName = "FS-ZHXC-3L-HH08-05-X-S9850-04"
	req.Vrid = 5
	req.Detail.IpAddr = "112.45.34.0"
	req.Detail.IpMask = "27"
	req.Detail.NeighborIp = "172.18.16.0"
	req.Detail.NeighborMask = "29"
	req.Detail.Operators = "CD-CMCC"
	ins, err := serviceImpl.VrrpDeployment(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestDoubleStaticDeployment(t *testing.T) {
	req := internet.NewDeploymentDoubleStatic()
	req.FirstDevName = "FS-ZHXC-3L-HH08-04-X-S9850-03"
	req.SecondDevName = "FS-ZHXC-3L-HH08-05-X-S9850-04"
	req.Detail.Operators = "CD-CT"
	req.Detail.IpAddr = "112.45.33.0"
	req.Detail.IpMask = "27"
	req.Detail.NeighborIp = "172.18.15.0"
	req.Detail.NeighborMask = "29"
	ins, err := serviceImpl.DoubleStaticDeployment(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestSingleDeployment(t *testing.T) {
	req := internet.NewDeploymentSingle()
	req.DevName = "FS-ZHXC-3L-HH08-04-X-S9850-03"
	req.Detail.Operators = "CD-CMCC"
	req.Detail.IpAddr = "112.45.29.0"
	req.Detail.IpMask = "27"
	req.Detail.NeighborIp = "172.18.15.0"
	req.Detail.NeighborMask = "30"
	ins, err := serviceImpl.SingleDeployment(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
