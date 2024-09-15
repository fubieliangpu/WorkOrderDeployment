package ioc_test

import (
	"testing"

	rcd "github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice/impl"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
)

func TestRegistry(t *testing.T) {
	ioc.Controller.Registry("rcdevice", &rcd.DeviceServiceImpl{})
	t.Logf("%p", ioc.Controller.Get("rcdevice"))
}
