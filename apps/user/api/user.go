package api

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
	"github.com/fubieliangpu/WorkOrderDeployment/exception"
	"github.com/fubieliangpu/WorkOrderDeployment/middleware"
	"github.com/fubieliangpu/WorkOrderDeployment/response"
	"github.com/gin-gonic/gin"
)

func (h *UserApiHandler) Registry(appRouter gin.IRouter) {
	appRouter.Use(middleware.Auth)
	appRouter.POST("/", middleware.RequireRole(user.ROLE_ADMIN), h.CreateUser)
	appRouter.PATCH("/", middleware.RequireRole(user.ROLE_ADMIN), h.ChangeUser)
	appRouter.DELETE("/", middleware.RequireRole(user.ROLE_ADMIN), h.DeleteUser)
}

// 创建用户 POST /wod/api/v1/users
func (h *UserApiHandler) CreateUser(ctx *gin.Context) {
	req := user.NewCreateUserRuquest()
	if err := ctx.Bind(req); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	ins, err := h.userapi.CreateUser(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}
	response.Success(ins, ctx)

}

// 修改用户密码 PATCH  /wod/api/v1/users
func (h *UserApiHandler) ChangeUser(ctx *gin.Context) {
	req := user.NewChangeUserRequest()

	//body
	if err := ctx.Bind(req); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	//业务处理，需用中间件处理鉴权及校验
	ins, err := h.userapi.ChangeUser(ctx, req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}
	response.Success(ins, ctx)
}

// 删除用户 DELETE /wod/api/v1/users
func (h *UserApiHandler) DeleteUser(ctx *gin.Context) {
	req := user.NewDeleteUserRequest()
	if err := ctx.Bind(req); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	//DeleteUser需用中间件处理鉴权及校验
	ins, err := h.userapi.DeleteUser(ctx, req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}
	response.Success(ins, ctx)
}
