package dci_test

import (
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/dci"
)

func TestNewCreateDCIRequest(t *testing.T) {
	req := dci.NewCreateDCIRequest()
	t.Log(req)
}

func TestNewDCI(t *testing.T) {
	req := dci.NewDCI()
	t.Log(req)
}
