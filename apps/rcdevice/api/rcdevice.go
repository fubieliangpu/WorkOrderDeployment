package api

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/token"
	"github.com/fubieliangpu/WorkOrderDeployment/common"
	"github.com/fubieliangpu/WorkOrderDeployment/exception"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"github.com/fubieliangpu/WorkOrderDeployment/response"
	"github.com/gin-gonic/gin"
)

func (h *RcDeviceApiHandler) Registry(appRouter gin.IRouter) {
	// 不需要鉴权，公开访问

	// 修改变更需要认证

}

// 设备条目列表 GET /wod/api/v1/rcdevice?page_size=10&page_number=1&...
func (h *RcDeviceApiHandler) QueryDeviceList(ctx *gin.Context) {
	//获取用户请求
	req := rcdevice.NewQueryDeviceListRequest()
	req.PageRequest = common.NewPageRequestFromGinCtx(ctx)
	req.KeyWords = ctx.Query("keywords")

	//业务处理
	set, err := h.svc.QueryDeviceList(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}

	//返回结果
	response.Success(set, ctx)
}

// 设备条目详情 GET /wod/api/v1/rcdevice/{device_name}
func (h *RcDeviceApiHandler) DescribeDevice(ctx *gin.Context) {
	tk, _ := ctx.Cookie(token.COOKIE_TOKEY_KEY)
	ioc.Controller.Get(token.AppName).(token.Service).ValidateToken(ctx.Request.Context(), token.NewValidateTokenRequest(tk))
	//获取用户请求
	req := rcdevice.NewDescribeDeviceRequest(ctx.Param("device_name"))
	//业务处理
	ins, err := h.svc.DescribeDevice(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}
	response.Success(ins, ctx)
}

// 创建设备条目 CreateDevice：POST /wod/api/v1/rcdevice
func (h *RcDeviceApiHandler) CreateDevice(ctx *gin.Context) {
	//获取用户请求
	req := rcdevice.NewCreateDeviceRequest()
	if err := ctx.Bind(req); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	//业务处理
	ins, err := h.svc.CreateDevice(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}
	//返回结果
	response.Success(ins, ctx)
}

// 设备全量修改接口 PUT /wod/api/v1/rcdevice/{device_name}
func (h *RcDeviceApiHandler) PutUpdateDevice(ctx *gin.Context) {
	//获取用户请求
	req := rcdevice.NewUpdateDeviceRequest(ctx.Param("device_name"))
	req.UpdateMode = common.UPDATE_MODE_PUT
	if err := ctx.Bind(req.CreateDeviceRequest); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	//业务处理
	ins, err := h.svc.UpdateDevice(ctx, req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}
	response.Success(ins, ctx)
}

// 设备增量修改接口 PATCH /wod/api/v1/rcdevice/{device_name}
func (h *RcDeviceApiHandler) PatchUpdateDevice(ctx *gin.Context) {
	//获取用户请求
	req := rcdevice.NewUpdateDeviceRequest(ctx.Param("device_name"))
	req.UpdateMode = common.UPDATE_MODE_PATCH

	//body parse and save in struct
	if err := ctx.Bind(req.CreateDeviceRequest); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}

	//业务处理
	ins, err := h.svc.UpdateDevice(ctx, req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}
	response.Success(ins, ctx)
}

// 设备删除接口 DELETE /wod/api/v1/rcdevice/{device_name}
func (h *RcDeviceApiHandler) DeleteDevice(ctx *gin.Context) {
	//获取用户请求
	req := rcdevice.NewDeleteDeviceRequest(ctx.Param("device_name"))

	//业务处理
	ins, err := h.svc.DeleteDevice(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}
	//返回结果
	response.Success(ins, ctx)
}

// 设备配置修改(相关状态查询) POST /wod/api/v1/rcdevice/{device_name}/config
func (h *RcDeviceApiHandler) ChangeDeviceConfig(ctx *gin.Context)  {
	//获取用户请求
	req := rcdevice.NewChangeDeviceConfigRequest(ctx.Param("device_name"),ctx.Param("device_config_file"),ctx.Param("user_file"))
	//body
	if err := ctx.Bind(req)
}
