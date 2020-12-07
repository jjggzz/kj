package discovery

type Discover interface {
	// 服务注册
	RegisterServer()
	// 服务注销
	DeregisterServer()
}
