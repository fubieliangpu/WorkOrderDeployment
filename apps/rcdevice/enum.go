package rcdevice

//设备状态已创建/已修改/已删除
type Status int

const (
	//已创建
	STATUS_CREATED = iota
	//已删除
	STATUS_DELETED
	//已修改
	STATUS_MODIFIED
)
