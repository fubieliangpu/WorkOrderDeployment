package internet_test

import (
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/internet"
)

func TestIpMaskOpt(t *testing.T) {
	req := internet.NewIpMaskOpt()
	t.Log(req)
}

func TestDeploymentVRRP(t *testing.T) {
	req := internet.NewDeploymentVRRP()
	t.Log(req)
}

func TestDeploymentDoubleStatic(t *testing.T) {
	req := internet.NewDeploymentDoubleStatic()
	t.Log(req)
}
