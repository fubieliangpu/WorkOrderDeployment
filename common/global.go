package common

import "time"

// 设备通用参数
type DevMeta struct {
	// 设备Id
	Id int `json:"id" gorm:"column:id"` //如果主键不为id(int)(之前主键用了name varchar(255))，updatedevice时会报WHERE conditions required异常
}

func NewDevMeta() *DevMeta {
	return &DevMeta{}
}

type Brand int

const (
	Cisco Brand = iota + 1
	Ruijie
	H3C
	Huawei
	Huawei_CE
	Huawei_FW
	Juniper
)

func NewUserMeta() *UserMeta {
	return &UserMeta{
		CreateAt: time.Now().Unix(),
	}
}

//用户通用参数
type UserMeta struct {
	//用户ID
	Id int `json:"id" gorm:"column:id"`
	//创建时间，时间戳 10位，秒
	CreateAt int64 `json:"created_at" gorm:"column:created_at"`
	//更新时间，时间戳 10位，秒
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at"`
}

type DeviceLevel int

const (
	CORE DeviceLevel = iota + 1
	CONVERGE
	Access
)
