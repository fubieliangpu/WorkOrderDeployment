package ioc

//Controller 是一个Container，使用MapContainer实现
var Controller Container = &MapContainer{
	name:   "controller",
	storge: make(map[string]Object),
}

//Api 所有的对外接口对象都放这里

var Api Container = &MapContainer{
	name:   "api",
	storge: make(map[string]Object),
}
