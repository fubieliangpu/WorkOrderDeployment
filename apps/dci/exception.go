package dci

import (
	"net/http"

	"github.com/fubieliangpu/WorkOrderDeployment/exception"
)

var (
	ErrDCIVlanconflict = exception.NewApiException(50054, "vlan相关配置冲突").WithHttpCode(http.StatusConflict)
	ErrDCIVNIconflict  = exception.NewApiException(50059, "vni配置冲突").WithHttpCode(http.StatusConflict)
	ErrDCIVSIconflict  = exception.NewApiException(50060, "vsi名称冲突").WithHttpCode(http.StatusConflict)
	ErrDCIBRconflict   = exception.NewApiException(50061, "聚合端口号冲突").WithHttpCode(http.StatusConflict)
	ErrDCINoTun        = exception.NewApiException(50062, "Tunnel口不存在").WithHttpCode(http.StatusConflict)
)
