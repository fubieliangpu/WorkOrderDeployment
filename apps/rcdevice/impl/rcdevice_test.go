package impl_test

import (
	"fmt"
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/common"
)

func TestCreateDevice(t *testing.T) {
	req := rcdevice.NewCreateDeviceRequest()
	req.Name = "FS-9850-01"
	req.ServerAddr = "192.168.11.229"
	req.Port = "22"
	req.Brand = 2
	req.Idc = "HYC"
	ins, err := serviceImpl.CreateDevice(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestDescribeDevice(t *testing.T) {
	req := rcdevice.NewDescribeDeviceRequest("RJJD-9850")
	ins, err := serviceImpl.DescribeDevice(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ins)
}

func TestQueryDeviceList(t *testing.T) {
	req := rcdevice.NewQueryDeviceListRequest()
	req.PageSize = 4
	req.PageNumber = 1
	ins, err := serviceImpl.QueryDeviceList(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ins)
}

func TestPatchUpdateDevice(t *testing.T) {
	req := rcdevice.NewUpdateDeviceRequest("FS-9850-01")
	req.UpdateMode = common.UPDATE_MODE_PATCH
	req.Brand = 0
	ins, err := serviceImpl.UpdateDevice(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestPutUpdateDevice(t *testing.T) {
	req := rcdevice.NewUpdateDeviceRequest("test")
	req.UpdateMode = common.UPDATE_MODE_PUT
	req.Name = "test"
	req.Brand = common.Huawei_FW
	req.ServerAddr = "192.168.79.10"
	req.Port = "22"
	req.Idc = "test"
	ins, err := serviceImpl.UpdateDevice(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestDeleteDevice(t *testing.T) {
	req := rcdevice.NewDeleteDeviceRequest("test")
	ins, err := serviceImpl.DeleteDevice(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestLoadPasswordFromYAML(t *testing.T) {
	req := rcdevice.NewDeviceUserInfo()
	ins, err := rcdevice.LoadUsernmPasswdFromYaml("./user.yml", req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestChangeDeviceConfig(t *testing.T) {
	req := rcdevice.NewChangeDeviceConfigRequest()
	req.DeviceName = "HYC-6890-01"
	req.DeviceConfigFile = "h3c-config.txt"
	req.UserFile = "user.yml"
	ins, err := serviceImpl.ChangeDeviceConfig(ctx, req)
	fmt.Println(ins.ServerAddr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
