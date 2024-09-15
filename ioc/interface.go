package ioc

//需要实现自动加载注册模块的初始化变量，只要实现了Init()方法就是Object对象
type Object interface {
	Init() error
}

//实现IOC托儿所
type Container interface {
	Registry(name string, obj Object)
	Get(name string) any
	Init() error
}
