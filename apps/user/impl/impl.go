package impl

import (
	"context"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/token"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
	"github.com/fubieliangpu/WorkOrderDeployment/common"
	"github.com/fubieliangpu/WorkOrderDeployment/conf"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"gorm.io/gorm"
)

//注册到ioc

func init() {
	ioc.Controller.Registry(user.AppName, &UserServiceImpl{})
}

type UserServiceImpl struct {
	db *gorm.DB
	tk token.Service
}

// 实现Init之后才能满足ioc.Object的要求
func (i *UserServiceImpl) Init() error {
	i.db = conf.C().MySQL.GetDB()
	i.tk = ioc.Controller.Get(token.AppName).(token.Service)
	return nil
}

func (i *UserServiceImpl) CreateUser(ctx context.Context, in *user.CreateUserRequest) (*user.User, error) {
	//校验请求
	if err := common.Validate(in); err != nil {
		return nil, err
	}
	//明文密码处理
	if err := in.HashPassword(); err != nil {
		return nil, err
	}
	//创建user对象
	ins := user.NewUser(in)

	//user持久化
	if err := i.db.WithContext(ctx).Save(ins).Error; err != nil {
		return nil, err
	}
	//返回存储了什么用户
	return ins, nil
}

// WHERE以及LIMIT条件语句查询
func (i *UserServiceImpl) QueryUser(ctx context.Context, in *user.QueryUserRequest) (*user.UserSet, error) {
	set := user.NewUserSet()
	//构造一个查询语句
	query := i.db.Model(&user.User{}).WithContext(ctx)
	if in.Username != "" {
		query = query.Where("username = ?", in.Username)
	}

	// 怎么查询Total, 需要把过滤条件: username ,key
	// 查询Total时能不能把分页参数带上
	// select COUNT(*) from xxx limit 10
	// select COUNT(*) from xxx
	// 不能携带分页参数
	if err := query.Count(&set.Total).Error; err != nil {
		return nil, err
	}

	if err := query.Offset(in.Offset()).Limit(in.PageSize).Find(&set.Items).Error; err != nil {
		return nil, err
	}
	return set, nil
}

func (i *UserServiceImpl) DeleteUser(ctx context.Context, in *user.DeleteUserRequest) (*user.User, error) {
	//首先校验Token的有效性
	//tkreq := token.NewValidateTokenRequest(in.AccessToken)
	// tkins, err := i.tk.ValidateToken(ctx, tkreq)
	// if err != nil {
	// 	return nil, err
	// }
	//校验用户身份是不是管理员
	// usqreq := user.NewQueryUserRequest()
	// usqreq.Username = in.Username
	// uset, err := i.QueryUser(ctx, usqreq)
	// if err != nil {
	// 	return nil, err
	// }
	//权限不为管理员则报错退出
	// if uset.Items[0].Role != user.ROLE_ADMIN {
	// 	return nil, user.ErrPermissionDeny
	// }

	//判定下所删除用户名是否是自己，如果是自己就不能删除，毕竟自己不能把自己举起来
	// if uset.Items[0].Username == in.Username {
	// 	return nil, user.ErrSameUsername
	// }
	//查询被删除用户
	req := user.NewQueryUserRequest()
	req.Username = in.Username
	ruset, err := i.QueryUser(ctx, req)
	if err != nil {
		return nil, err
	}
	//先删除令牌，再删除用户
	err = i.db.WithContext(ctx).Table("tokens").Where("username = ?", in.Username).Delete(token.Token{}).Error
	if err != nil {
		return nil, err
	}
	err = i.db.WithContext(ctx).Table("users").Where("username = ?", in.Username).Delete(user.User{}).Error
	if err != nil {
		return nil, err
	}
	//返回删了什么用户
	return ruset.Items[0], nil
}

// 修改用户密码
func (i *UserServiceImpl) ChangeUser(ctx context.Context, in *user.ChangeUserRequest) (*user.User, error) {
	//首先校验Token的有效性
	// tkreq := token.NewValidateTokenRequest(in.AccessToken)
	// tkins, err := i.tk.ValidateToken(ctx, tkreq)
	// if err != nil {
	// 	return nil, err
	// }
	//判断用户名是否相同，即普通用户也可以修改自己身密码
	//校验用户身份是不是管理员
	// usqreq := user.NewQueryUserRequest()
	// usqreq.Username = in.Username
	// uset, err := i.QueryUser(ctx, usqreq)
	// if err != nil {
	// 	return nil, err
	// }
	// if in.Username == tkins.UserName {
	// 	uset.Items[0].Password = in.Password
	// 	err = uset.Items[0].HashPassword()
	// 	//持久化
	// 	if err := i.db.WithContext(ctx).Save(uset.Items[0]).Error; err != nil {
	// 		return nil, err
	// 	}
	// 	return uset.Items[0], err
	// }
	// //权限不为管理员则报错退出
	// //如果权限是管理员，则修改密码

	// usqreq.Username = tkins.UserName
	// //查询tk对应user是否为管理员
	// aduset, err := i.QueryUser(ctx, usqreq)
	// if aduset.Items[0].Role != user.ROLE_ADMIN {
	// 	return nil, user.ErrPermissionDeny
	// } else if aduset.Items[0].Role == user.ROLE_ADMIN {
	// 	uset.Items[0].Password = in.Password
	// 	err = uset.Items[0].HashPassword()
	// 	if err := i.db.WithContext(ctx).Save(uset.Items[0]).Error; err != nil {
	// 		return nil, err
	// 	}
	// 	return uset.Items[0], err
	// } else {
	// 	return nil, err
	// }
	//以上鉴权判定交给中间件
	req := user.NewQueryUserRequest()
	req.Username = in.Username
	ruset, err := i.QueryUser(ctx, req)
	if err != nil {
		return nil, err
	}
	if ruset.Total == 0 {
		return nil, user.ErrNoUser
	}
	ruset.Items[0].Password = in.Password
	err = ruset.Items[0].HashPassword()
	if err := i.db.WithContext(ctx).Save(ruset.Items[0]).Error; err != nil {
		return nil, err
	}
	return ruset.Items[0], err
}
