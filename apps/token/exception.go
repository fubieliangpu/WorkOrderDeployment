package token

import (
	"net/http"

	"github.com/fubieliangpu/WorkOrderDeployment/exception"
)

var (
	ErrAccessTokenExpired  = exception.NewApiException(50002, "AccessToken过期")
	ErrRefreshTokenExpired = exception.NewApiException(50003, "RefreshToken过期")
	ErrAuthFailed          = exception.NewApiException(50001, "用户名或者密码不正确").WithHttpCode(http.StatusUnauthorized)
	ErrUnauthorized        = exception.NewApiException(50400, "请登录").WithHttpCode(http.StatusUnauthorized)
	ErrPermissionDeny      = exception.NewApiException(50004, "当前角色无权限访问该接口").WithHttpCode(http.StatusForbidden)
)
