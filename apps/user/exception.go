package user

import (
	"net/http"

	"github.com/fubieliangpu/WorkOrderDeployment/exception"
)

var (
	ErrPermissionDeny = exception.NewApiException(50044, "用户无权限").WithHttpCode(http.StatusForbidden)
	ErrSameUsername   = exception.NewApiException(50043, "用户名相同不能删除").WithHttpCode(http.StatusForbidden)
	ErrNoUser         = exception.NewApiException(50045, "查询不到此用户").WithHttpCode(http.StatusNotFound)
)
