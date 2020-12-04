package discovery

import "google.golang.org/grpc"

type Discover interface {
	// 服务注册
	RegisterServer() error
	// 服务发现
	DiscoveryServers([]string) error
	// 根据服务名获取服务连接列表
	GetConn(string) ([]*grpc.ClientConn, error)
	//获取一个服务连接
	GetOneConn(string) (*grpc.ClientConn, error)
}

type Instance struct {
	Address string
	Conn    *grpc.ClientConn
}
