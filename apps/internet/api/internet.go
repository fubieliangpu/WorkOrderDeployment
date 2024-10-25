package api

import (
	"fmt"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/internet"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
	"github.com/fubieliangpu/WorkOrderDeployment/exception"
	"github.com/fubieliangpu/WorkOrderDeployment/middleware"
	"github.com/fubieliangpu/WorkOrderDeployment/notice/byemail"
	"github.com/fubieliangpu/WorkOrderDeployment/response"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
)

func (h *InternetApiHandler) Registry(appRouter gin.IRouter) {
	//修改变更需要认证
	appRouter.Use(middleware.Auth)
	appRouter.GET("/vrrpcheck", middleware.RequireRole(user.ROLE_ADMIN), h.VrrpConflictCheck)
	appRouter.GET("/doublecheck", middleware.RequireRole(user.ROLE_ADMIN), h.DoubleStaticConflictCheck)
	appRouter.GET("/singlecheck", middleware.RequireRole(user.ROLE_ADMIN), h.SingleConflictCheck)
	appRouter.PUT("/vrrpdeployment", middleware.RequireRole(user.ROLE_ADMIN), h.VrrpDeployment)
	appRouter.PUT("/doubledeployment", middleware.RequireRole(user.ROLE_ADMIN), h.DoubleStaticDeployment)
	appRouter.PUT("/singledeployment", middleware.RequireRole(user.ROLE_ADMIN), h.SingleDeployment)
}

// 部署vrrp前冲突检测 /wod/api/v1/internet/vrrpcheck
func (h *InternetApiHandler) VrrpConflictCheck(ctx *gin.Context) {
	//获取用户请求
	req := internet.NewDeploymentVRRP()
	if err := ctx.BindJSON(req); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	//业务处理
	status, err := h.svc.VrrpConflictCheck(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}
	//返回结果
	response.Success(status, ctx)
}

// 部署双上联前冲突检测 /wod/api/v1/internet/doublecheck
func (h *InternetApiHandler) DoubleStaticConflictCheck(ctx *gin.Context) {
	//获取用户请求
	req := internet.NewDeploymentDoubleStatic()
	if err := ctx.BindJSON(req); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	//业务处理
	status, err := h.svc.DoubleStaticConflictCheck(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}
	//返回结果
	response.Success(status, ctx)
}

// 部署单上联前冲突检测 /wod/api/v1/internet/singlecheck
func (h *InternetApiHandler) SingleConflictCheck(ctx *gin.Context) {
	//获取用户请求
	req := internet.NewDeploymentSingle()
	if err := ctx.BindJSON(req); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	//业务处理
	status, err := h.svc.SingleConflictCheck(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(err, ctx)
		return
	}
	//返回结果
	response.Success(status, ctx)
}

// 部署业务/wod/api/v1/internet/vrrpdeployment
func (h *InternetApiHandler) VrrpDeployment(ctx *gin.Context) {
	//获取用户请求
	req := internet.NewDeploymentVRRP()
	if err := ctx.BindJSON(req); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	//业务处理

	result, err := h.svc.VrrpDeployment(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(err, ctx)
	}
	//构建邮件认证信息
	mauthinfo := byemail.NewMailAuthInfo()
	if _, err := mauthinfo.LoadEmailAuthInfoFromYaml("mailauth.yaml"); err != nil {
		response.Failed(err, ctx)
	}
	//构建邮件发送的内容
	msg := fmt.Sprintf(
		"VRRP下发结果\n运营商: %v\n业务IP: %v/%v\n互联IP段: %v/%v\n下发设备: %v,%v\n配置下发完成,请登陆设备检查结果!\n",
		req.Detail.Operators,
		req.Detail.IpAddr,
		req.Detail.IpMask,
		req.Detail.NeighborIp,
		req.Detail.NeighborMask,
		req.MasterDevName,
		req.BackupDevName,
	)
	//构建邮件对象
	e := email.NewEmail()
	e.From = fmt.Sprintf("WorkOrderNotice <%v>", mauthinfo.SendAddress)
	e.Subject = "业务下发结果"
	e.To = []string{mauthinfo.RecvAddress}
	e.Text = []byte(msg)

	//发送邮件
	if err := mauthinfo.MySendmail(e); err != nil {
		response.Failed(err, ctx)
	}
	//返回结果
	response.Success(result, ctx)
}

// 部署业务/wod/api/v1/internet/doubledeployment
func (h *InternetApiHandler) DoubleStaticDeployment(ctx *gin.Context) {
	//获取用户请求
	req := internet.NewDeploymentDoubleStatic()
	if err := ctx.BindJSON(req); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	//业务处理

	result, err := h.svc.DoubleStaticDeployment(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(err, ctx)
	}

	//构建邮件认证信息
	mauthinfo := byemail.NewMailAuthInfo()
	if _, err := mauthinfo.LoadEmailAuthInfoFromYaml("mailauth.yaml"); err != nil {
		response.Failed(err, ctx)
	}
	//构建邮件发送的内容
	msg := fmt.Sprintf(
		"运营商：%v\n业务IP: %v/%v\n互联IP段: %v/%v\n下发设备: %v,%v\n配置下发完成,请登陆设备检查结果!\n",
		req.Detail.Operators,
		req.Detail.IpAddr,
		req.Detail.IpMask,
		req.Detail.NeighborIp,
		req.Detail.NeighborMask,
		req.FirstDevName,
		req.SecondDevName,
	)
	//构建邮件对象
	e := email.NewEmail()
	e.From = fmt.Sprintf("WorkOrderNotice <%v>", mauthinfo.SendAddress)
	e.Subject = "业务下发结果"
	e.To = []string{mauthinfo.RecvAddress}
	e.Text = []byte(msg)
	//发送邮件
	if err := mauthinfo.MySendmail(e); err != nil {
		response.Failed(err, ctx)
	}
	//返回结果
	response.Success(result, ctx)
}

// 部署业务/wod/api/v1/internet/singledeployment
func (h *InternetApiHandler) SingleDeployment(ctx *gin.Context) {
	//获取用户请求
	req := internet.NewDeploymentSingle()
	if err := ctx.BindJSON(req); err != nil {
		response.Failed(exception.ErrValidateFailed(err.Error()), ctx)
		return
	}
	//业务处理

	result, err := h.svc.SingleDeployment(ctx.Request.Context(), req)
	if err != nil {
		response.Failed(err, ctx)
	}

	//构建邮件认证信息
	mauthinfo := byemail.NewMailAuthInfo()
	if _, err := mauthinfo.LoadEmailAuthInfoFromYaml("mailauth.yaml"); err != nil {
		response.Failed(err, ctx)
	}
	//构建邮件发送的内容
	msg := fmt.Sprintf(
		"运营商：%v\n业务IP: %v/%v\n互联IP段: %v/%v\n下发设备: %v\n配置下发完成,请登陆设备检查结果!\n",
		req.Detail.Operators,
		req.Detail.IpAddr,
		req.Detail.IpMask,
		req.Detail.NeighborIp,
		req.Detail.NeighborMask,
		req.DevName,
	)
	//构建邮件对象
	e := email.NewEmail()
	e.From = fmt.Sprintf("WorkOrderNotice <%v>", mauthinfo.SendAddress)
	e.Subject = "业务下发结果"
	e.To = []string{mauthinfo.RecvAddress}
	e.Text = []byte(msg)
	//发送邮件
	if err := mauthinfo.MySendmail(e); err != nil {
		response.Failed(err, ctx)
	}

	//返回结果
	response.Success(result, ctx)
}
