package impl

import (
	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"gorm.io/gorm"
)

func init() {
	ioc.Controller.Registry(rcdevice.AppName, &DeviceServiceImpl{})
}

type DeviceServiceImpl struct {
	db *gorm.DB
}

func (i *DeviceServiceImpl) Init() error {
	return nil
}
