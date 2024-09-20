package impl

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
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
