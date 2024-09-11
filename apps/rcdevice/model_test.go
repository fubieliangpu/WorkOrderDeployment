package rcdevice_test

import (
	"fmt"
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
)

func TestNewDevice(t *testing.T) {
	ins := rcdevice.NewDevice()
	t.Log(ins.CreateDeviceRequest)
}

func TestNewChangedDeviceStatusRequest(t *testing.T) {
	ins := rcdevice.NewChangedDeviceStatusRequest()
	ins.SetStatus(rcdevice.STATUS_CREATED)
	t.Log(ins)
}

func TestNewDeviceSet(t *testing.T) {
	idclist := [3]string{"rjjd", "hyc", "ly"}
	ins := rcdevice.NewDeviceSet()
	for i := 0; i < len(idclist); i++ {
		ins.Items = append(ins.Items, rcdevice.NewDevice().SetIDC(idclist[i]))
		//fmt.Printf("%T,%[1]v\n", rcdevice.NewDevice().SetIDC(idclist[i]))
	}
	ins.Total = int64(len(ins.Items))
	fmt.Printf("%+v\n", ins)
}
