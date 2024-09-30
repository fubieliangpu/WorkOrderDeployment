package internet

import (
	"net/http"

	"github.com/fubieliangpu/WorkOrderDeployment/exception"
)

var (
	ErrNoDeviceInIdc   = exception.NewApiException(50048, "指定IDC查询不到设备").WithHttpCode(http.StatusNotFound)
	ErrRouteConflict   = exception.NewApiException(50050, "所下发路由冲突").WithHttpCode(http.StatusConflict)
	ErrBrandNotSupport = exception.NewApiException(50055, "不支持该品牌设备").WithHttpCode(http.StatusNotAcceptable)
)
