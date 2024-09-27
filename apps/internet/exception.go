package internet

import (
	"net/http"

	"github.com/fubieliangpu/WorkOrderDeployment/exception"
)

var (
	ErrNoDeviceInIdc = exception.NewApiException(50048, "指定IDC查询不到设备").WithHttpCode(http.StatusNotFound)
)
