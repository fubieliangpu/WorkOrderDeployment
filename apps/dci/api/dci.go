package api

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/dci"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
	"github.com/fubieliangpu/WorkOrderDeployment/exception"
	"github.com/fubieliangpu/WorkOrderDeployment/middleware"
	"github.com/fubieliangpu/WorkOrderDeployment/response"
	"github.com/gin-gonic/gin"
)

func (h *DCIApiHandler) Registry(appRouter gin.IRouter) {
	//修改变更需要认证
	appRouter.Use(middleware.Auth)
	appRouter.GET("/conflictcheck", middleware.RequireRole(user.ROLE_ADMIN), h.ConflictCheck)
	appRouter.PUT("/deploymentdci", middleware.RequireRole(user.ROLE_ADMIN), h.ConfigDeployment)
	appRouter.DELETE("/revokedeployment", middleware.RequireRole(user.ROLE_ADMIN), h.ConfigDeployment)
}

// /wod/api/v1/dci/conflictcheck
func (h *DCIApiHandler) ConflictCheck(ctx *gin.Context) {
	//获取用户请求
	req := dci.NewCreateDCIRequest()
	if err := ctx.BindJSON(req); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	//业务处理
	if err := h.svc.ConflictCheck(ctx.Request.Context(), req); err != nil {
		response.Failed(err, ctx)
		return
	}

	//检查远端设备
	req.SrcDevName, req.DestDevName = req.DestDevName, req.SrcDevName
	if err := h.svc.ConflictCheck(ctx.Request.Context(), req); err != nil {
		response.Failed(err, ctx)
		return
	}

	//返回结果
	response.Success("Check no conflict", ctx)
}

// /wod/api/v1/dci/deploymentdci
func (h *DCIApiHandler) ConfigDeployment(ctx *gin.Context) {
	//获取用户请求
	req := dci.NewCreateDCIRequest()
	if err := ctx.BindJSON(req); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	//业务处理
	resdci, err := h.svc.ConfigDeployment(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}
	//返回结果
	response.Success(resdci, ctx)
}
