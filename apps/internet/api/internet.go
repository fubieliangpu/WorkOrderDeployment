package api

// func (h *InternetApiHandler) Registry(appRouter gin.IRouter) {
// 	//修改变更需要认证
// 	appRouter.Use(middleware.Auth)
// 	appRouter.GET("/conflictcheck", middleware.RequireRole(user.ROLE_ADMIN), h.ConflictCheck)
// 	appRouter.PUT("/deploymentnetpd", middleware.RequireRole(user.ROLE_ADMIN), h.ConfigDeployment)
// 	appRouter.DELETE("/revokedeployment", middleware.RequireRole(user.ROLE_ADMIN), h.RevokeDeployment)
// }

// 部署业务前冲突检测 /wod/api/v1/internet/conflictcheck
// func (h *InternetApiHandler) ConflictCheck(ctx *gin.Context) {
// 	//获取用户请求
// 	req := internet.NewDeploymentNetworkProductRequest()
// 	if err := ctx.BindJSON(req); err != nil {
// 		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
// 		return
// 	}
// 	//业务处理
// 	status, err := h.svc.ConflictCheck(ctx.Request.Context(), req)
// 	if err != nil {
// 		response.Failed(err, ctx)
// 		return
// 	}
// 	//返回结果
// 	response.Success(status, ctx)
// }

// 部署业务/wod/api/v1/internet/deploymentnetpd
// func (h *InternetApiHandler) ConfigDeployment(ctx *gin.Context) {
// 	//获取用户请求
// 	req := internet.NewDeploymentNetworkProductRequest()
// 	if err := ctx.BindJSON(req); err != nil {
// 		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
// 		return
// 	}
// 	//业务处理
// 	if req.ConfigRevoke == 0 {
// 		netp, err := h.svc.ConfigDeployment(ctx.Request.Context(), req)
// 		if err != nil {
// 			response.Failed(err, ctx)
// 			return
// 		}
// 		//返回结果
// 		response.Success(netp, ctx)
// 	} else {
// 		response.Failed(exception.ErrValidateFailed("非业务部署,API调用错误"), ctx)
// 		return
// 	}

// }

// 配置回收/wod/api/v1/internet/revokedeployment
// func (h *InternetApiHandler) RevokeDeployment(ctx *gin.Context) {
// 	//获取用户请求
// 	req := internet.NewDeploymentNetworkProductRequest()
// 	if err := ctx.BindJSON(req); err != nil {
// 		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
// 		return
// 	}
// 	//业务处理
// 	if req.ConfigRevoke == 1 {
// 		netp, err := h.svc.ConfigDeployment(ctx.Request.Context(), req)
// 		if err != nil {
// 			response.Failed(err, ctx)
// 			return
// 		}
// 		//返回结果
// 		response.Success(netp, ctx)
// 	} else {
// 		response.Failed(exception.ErrValidateFailed("非业务撤销,API调用错误"), ctx)
// 		return
// 	}
// }
