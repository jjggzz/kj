package discovery

import "github.com/go-kit/kit/sd"

type Discover interface {
	// 服务注册
	RegisterServer()
	// 服务注销
	DeregisterServer()
	// 发现服务
	Discovery(targetServerName string) (sd.Instancer, error)
}
