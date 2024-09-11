package rcdevice

//设备状态创建/修改/删除
type Status int

const (
	//创建
	STATUS_CREATED = iota
	//删除
	STATUS_DELETED
	//修改
	STATUS_MODIFIED
)
