package api

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/token"
	"github.com/fubieliangpu/WorkOrderDeployment/conf"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"github.com/fubieliangpu/WorkOrderDeployment/response"
	"github.com/gin-gonic/gin"
)

// 注册到ioc中
func init() {
	ioc.Api.Registry(token.AppName, &TokenApiHandler{})
}

// API Handler
type TokenApiHandler struct {
	//依赖Token接口，跟rcdevice类似
	token token.Service
}

func (h *TokenApiHandler) Registry(appRouter gin.IRouter) {
	//待完善
}

func (h *TokenApiHandler) Init() error {
	h.token = ioc.Controller.Get(token.AppName).(token.Service)
	//定义路由/wod/api/v1/tokens
	subRouter := conf.C().Application.GinRootRouter().Group("tokens")
	h.Registry(subRouter)
	return nil
}

// 登陆后颁发令牌
func (h *TokenApiHandler) Login(c *gin.Context) {
	//获取HTTP请求，颁发令牌请求
	req := token.NewIssueTokenRequest("", "")
	if err := c.BindJSON(req); err != nil {
		response.Failed(err, c)
		return
	}
	//业务处理
	tk, err := h.token.IssueToken(c.Request.Context(), req)
	if err != nil {
		response.Failed(err, c)
		return
	}
	//返回结果
	c.SetCookie(
		token.COOKIE_TOKEY_KEY,
		tk.AccessToken,
		tk.RefreshTokenExpiredAt,
		"/",
		conf.C().Application.Domain,
		false,
		true,
	)
	response.Success(tk, c)
}

// 退出登陆后撤销令牌
func (h *TokenApiHandler) Logout(c *gin.Context) {
	//DELETE方法一般不携带body
	ak, err := c.Cookie(token.COOKIE_TOKEY_KEY)
	if err != nil {
		response.Failed(err, c)
		return
	}
	rt := c.GetHeader(token.REFRESH_HEADER_KEY)
	req := token.NewRevolkTokenRequest(ak, rt)
	tk, err := h.token.RevolkToken(c.Request.Context(), req)
	if err != nil {
		response.Failed(err, c)
		return
	}
	response.Success(tk, c)
}
