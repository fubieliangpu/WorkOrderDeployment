package middleware

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/token"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
	"github.com/fubieliangpu/WorkOrderDeployment/response"
	"github.com/gin-gonic/gin"
)

// 这是一个需要有参数的中间件: Require("admin")
// 通过一个函数返回一个中间件: gin HandleFunc
// 这个中间件是加载在 认证中间件之后的
func RequireRole(requiredRoles ...user.Role) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		//判断当前用户的身份能不能匹配角色
		//补充上下文中注入的中间数据
		if v, ok := ctx.Get(token.GIN_TOKEN_KEY_NAME); ok {
			for i := range requiredRoles {
				requiredRole := requiredRoles[i]
				if v.(*token.Token).Role == requiredRole {
					ctx.Next()
					return
				}
			}
		}
		response.Failed(token.ErrPermissionDeny, ctx)
		ctx.Abort()
	}
}
