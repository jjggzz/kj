package discovery

import "google.golang.org/grpc"

// 服务注册的接口，所有服务都需继承此接口
type RpcServer interface{}

type Discovery interface {
	// 服务注册
	RegisterService([]RpcServer) error
	// 服务发现
	DiscoveryService([]RpcServer) error
	// 获取服务的连接
	GetConn(server RpcServer) (*grpc.ClientConn, error)
}

type Instance struct {
	// 服务地址
	Address string
	// grpc的连接
	Conn *grpc.ClientConn
}
