package middleware

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/token"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"github.com/fubieliangpu/WorkOrderDeployment/response"
	"github.com/gin-gonic/gin"
)

// Gin Web中间件,  我们需要在中间件注入到请求的链路当中，然后由Gin框架来调用
// HandlerFunc defines the handler used by gin middleware as return value.
// type HandlerFunc func(*Context)
// 加一下中间件处理函数 GIn HandlerFunc

func Auth(ctx *gin.Context) {
	accessToken, err := ctx.Cookie(token.COOKIE_TOKEY_KEY)
	if err != nil {
		response.Failed(token.ErrUnauthorized.WithMessage(err.Error()), ctx)
		ctx.Abort()
	}
	tk, err := ioc.Controller.Get(token.AppName).(token.Service).ValidateToken(ctx.Request.Context(), token.NewValidateTokenRequest(accessToken))
	if err != nil {
		response.Failed(token.ErrAuthFailed.WithMessage(err.Error()), ctx)
		ctx.Abort()
	} else {
		// 后面的handler 怎么知道 鉴权成功了, 当前是谁在访问这个接口
		// 请求的上下文:
		// 怎么把中间件请求结果，添加到请求的上下文中
		// 	// Keys is a key/value pair exclusively for the context of each request.
		// Keys map[string]any
		// Gin 采用一个map对象来维护中间传递的数据
		// context.WithValue()
		ctx.Set(token.GIN_TOKEN_KEY_NAME, tk)
		ctx.Next()
	}
}
