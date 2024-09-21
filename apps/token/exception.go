package token

import (
	"net/http"

	"github.com/fubieliangpu/WorkOrderDeployment/exception"
)

var (
	ErrAccessTokenExpired  = exception.NewApiException(50002, "AccessToken过期")
	ErrRefreshTokenExpired = exception.NewApiException(50003, "RefreshToken过期")
	ErrAuthFailed          = exception.NewApiException(50001, "用户名或者密码不正确").WithHttpCode(http.StatusUnauthorized)
)
