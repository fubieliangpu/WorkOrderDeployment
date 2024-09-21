package token

import "github.com/fubieliangpu/WorkOrderDeployment/exception"

var (
	ErrAccessTokenExpired  = exception.NewApiException(50002, "AccessToken过期")
	ErrRefreshTokenExpired = exception.NewApiException(50003, "RefreshToken过期")
)
