package common

// 通用参数
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
