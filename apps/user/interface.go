package user

import (
	"context"

	"github.com/fubieliangpu/WorkOrderDeployment/common"
)

//定义注册到ioc的Container的应用名

const (
	AppName = "user"
)

// user.Service
// 用户管理接口
// 接口定义的原则:  站在调用方(使用者)的角度来设计接口
// userServiceImpl.CreateUser(ctx, *CreateUserRequest)

type Service interface {
	//创建用户，只有管理员才能创建用户
	CreateUser(context.Context, *CreateUserRequest) (*User, error)
	//查询用户
	QueryUser(context.Context, *QueryUserRequest) (*UserSet, error)
	//删除用户,只有管理员才能删除用户
	DeleteUser(context.Context, *DeleteUserRequest) (*User, error)
	//修改用户密码，只有自己和管理员才能修改用户密码
	ChangeUser(context.Context, *ChangeUserRequest) (*User, error)
}

type QueryUserRequest struct {
	Username string
	*common.PageRequest
}

func NewQueryUserRequest() *QueryUserRequest {
	return &QueryUserRequest{
		PageRequest: common.NewPageRequest(),
	}
}

type DeleteUserRequest struct {
	//需要校验用户令牌及用户身份
	//需要被删除的用户名
	Username string `json:"username"`
	//以什么身份删除用户,通过中间件鉴权
	//AccessToken string `json:"access_token"`
}

func NewDeleteUserRequest() *DeleteUserRequest {
	return &DeleteUserRequest{}
}

type ChangeUserRequest struct {
	//需要被修改的用户的用户名及期望的密码
	Username string `json:"username"`
	Password string `json:"password"`
	//以什么身份修改用户,鉴权交给中间件
	//AccessToken string `json:"access_token"`
}

func NewChangeUserRequest() *ChangeUserRequest {
	return &ChangeUserRequest{}
}
