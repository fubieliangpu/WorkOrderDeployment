package impl

import (
	"context"

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
}

// 实现Init之后才能满足ioc.Object的要求
func (i *UserServiceImpl) Init() error {
	i.db = conf.C().MySQL.GetDB()
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
	return nil, nil
}

func (i *UserServiceImpl) ChangeUser(ctx context.Context, in *user.ChangeUserRequest) (*user.User, error) {
	return nil, nil
}
