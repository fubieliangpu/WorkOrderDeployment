package api

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
	"github.com/fubieliangpu/WorkOrderDeployment/middleware"
	"github.com/gin-gonic/gin"
)

func (h *InternetApiHandler) Registry(appRouter gin.IRouter) {
	//修改变更需要认证
	appRouter.GET("/conflictcheck", middleware.RequireRole(user.ROLE_ADMIN), h.ConflictCheck)

}

// 部署业务前冲突检测 /wod/api/v1/internet/conflictcheck
func (h *InternetApiHandler) ConflictCheck(ctx *gin.Context) {
	//获取用户请求
	//业务处理
	//返回结果
}
