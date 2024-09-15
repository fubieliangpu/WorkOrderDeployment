package ioc

//Controller 是一个Container，使用MapContainer实现
var Controller Container = &MapContainer{
	name:   "controller",
	storge: make(map[string]Object),
}

//
