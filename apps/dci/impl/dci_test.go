package impl_test

import (
	"log"
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/dci"
)

func TestConflictCheck(t *testing.T) {
	req := dci.NewCreateDCIRequest()
	req.BridgePort = "6"
	req.DestDevName = "SZ-HYC-1L-1-051-V-S6800-01"
	req.SrcDevName = "SZ-JYJF-110607-V-S9850-01"
	req.Idc = "JYJF"
	req.PortNumber = 1
	req.Vlan = 2223
	req.Vname = "HYC-01-SZJY-01-01-SZCT-2223"
	req.Vni = 2223
	err := serviceImpl.ConflictCheck(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(err)
}

func TestConfigDeployment(t *testing.T) {
	req := dci.NewCreateDCIRequest()
	req.BridgePort = "6"
	req.DestDevName = "SZ-HYC-1L-1-051-V-S6800-01"
	req.SrcDevName = "SZ-JYJF-110607-V-S9850-01"
	req.Idc = "JYJF"
	req.PortNumber = 1
	req.Vlan = 2223
	req.Vname = "HYC-01-SZJY-01-01-SZCT-2223"
	req.Vni = 2223
	mdci, err := serviceImpl.ConfigDeployment(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(mdci)

}
