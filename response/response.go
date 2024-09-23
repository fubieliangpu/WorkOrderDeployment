package response

import (
	"net/http"

	"github.com/fubieliangpu/WorkOrderDeployment/exception"
	"github.com/gin-gonic/gin"
)

func Success(data any, c *gin.Context) {
	c.JSON(http.StatusOK, data)
}

func Failed(err error, c *gin.Context) {
	//非200，返回ApiException对象
	httpCode := http.StatusInternalServerError
	if v, ok := err.(*exception.ApiException); ok {
		if v.HttpCode != 0 {
			httpCode = v.HttpCode
		}
	} else {
		//非业务异常，内部报错异常
		err = exception.ErrServerInternal(err.Error())
	}
	c.JSON(httpCode, err)
	c.Abort()
}
