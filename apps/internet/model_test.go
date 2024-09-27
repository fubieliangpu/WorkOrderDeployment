package internet_test

import (
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/internet"
)

func TestNewNetProd(t *testing.T) {
	req := internet.NewNetProd()
	t.Log(req)
}

func TestUndoConfig(t *testing.T) {
	req := internet.NewUndoDeviceConfigRequest()
	t.Log(req)
}
