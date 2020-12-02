package discovery

import "google.golang.org/grpc"

type Discover interface {
	// 服务注册
	RegisterServer() error
	// 服务发现
	DiscoveryServers([]string) error
	// 获取服务的连接
	GetConn(string) (*grpc.ClientConn, error)
}

type Instance struct {
	Address string
	Conn    *grpc.ClientConn
}
