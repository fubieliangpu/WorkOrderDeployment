package common

//设备通用参数
type DevMeta struct {
	//设备工作机房
	IDC string `json:"id" gorm:"column:idc"`
}
