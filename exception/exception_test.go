package exception_test

import (
	"encoding/json"
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/exception"
)

func GiveError() error {
	return exception.NewApiException(50001, "设备不存在")
}

func TestException(t *testing.T) {
	err := GiveError()

	t.Log(err)

	if v, ok := err.(*exception.ApiException); ok {
		t.Log(v.Code)
		t.Log(v.String())
	}

	dj, _ := json.MarshalIndent(err, "", "	")
	t.Log(string(dj))
}
