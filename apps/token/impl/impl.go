package impl

import (
	"context"
	"fmt"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/token"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
	"github.com/fubieliangpu/WorkOrderDeployment/conf"
	"github.com/fubieliangpu/WorkOrderDeployment/exception"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"gorm.io/gorm"
)

type TokenServiceImpl struct {
	db   *gorm.DB
	user user.Service
}

func init() {
	ioc.Controller.Registry(token.AppName, &TokenServiceImpl{})
}

// 实现ioc.Object
func (i *TokenServiceImpl) Init() error {
	i.db = conf.C().MySQL.GetDB()
	i.user = ioc.Controller.Get(user.AppName).(user.Service)
	return nil
}

// 令牌颁发
func (i *TokenServiceImpl) IssueToken(ctx context.Context, in *token.IssueTokenRequest) (*token.Token, error) {
	//查询用户对象
	queryUser := user.NewQueryUserRequest()
	queryUser.Username = in.Username
	us, err := i.user.QueryUser(ctx, queryUser)
	if err != nil {
		return nil, err
	}

	//匹配到这个用户时才能继续，否则不能继续
	if len(us.Items) == 0 {
		return nil, token.ErrAuthFailed
	}

	//比对用户名密码

	u := us.Items[0]
	if err := u.CheckPassword(in.Password); err != nil {
		return nil, token.ErrAuthFailed
	}

	//用户名密码认证没问题就颁发令牌
	tk := token.NewToken(u)

	//令牌入库
	if err := i.db.WithContext(ctx).Create(tk).Error; err != nil {
		return nil, exception.ErrServerInternal("保存令牌，%s", err)
	}
	//返回所保存的令牌
	return tk, nil
}

// 令牌撤销
func (i *TokenServiceImpl) RevolkToken(ctx context.Context, in *token.RevolkTokenRequest) (*token.Token, error) {
	//查询出Token
	tk := token.DefaultToken()
	err := i.db.WithContext(ctx).Where("access_token = ?", in.AccessToken).First(tk).Error

	if err == gorm.ErrRecordNotFound {
		return nil, exception.ErrNotFound("Token不存在")
	}

	if tk.RefreshToken != in.RefreshToken {
		return nil, fmt.Errorf("RefreshToken不正确")
	}
	//删除Token
	err = i.db.WithContext(ctx).Where("access_token = ?", in.AccessToken).Delete(token.Token{}).Error
	if err != nil {
		return nil, err
	}
	return tk, nil
}

// 令牌校验
func (i *TokenServiceImpl) ValidateToken(ctx context.Context, in *token.ValidateTokenRequest) (*token.Token, error) {
	tk := token.DefaultToken()
	//查询Token
	err := i.db.WithContext(ctx).Where("access_token = ?", in.AccessToken).First(tk).Error
	if err == gorm.ErrRecordNotFound {
		return nil, exception.ErrNotFound("Token不存在")
	}

	if err != nil {
		return nil, exception.ErrServerInternal("查询错误，%s", err)
	}

	//判断令牌有效性
	if err := tk.RefreshTokenIsExpired(); err != nil {
		return nil, err
	}
	if err := tk.AccessTokenIsExpired(); err != nil {
		return nil, err
	}

	//用户角色比较，成功后返回所校验的tk
	queryUserReq := user.NewQueryUserRequest()
	queryUserReq.Username = tk.UserName
	us, err := i.user.QueryUser(ctx, queryUserReq)
	if err != nil {
		return nil, err
	}
	if len(us.Items) == 0 {
		return nil, fmt.Errorf("token user not found")
	}
	tk.Role = us.Items[0].Role
	return tk, nil
}
