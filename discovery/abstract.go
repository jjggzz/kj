package discovery

type Discover interface {
	// 服务注册
	RegisterServer() error
	// 服务发现
	DiscoveryServers([]string) error
	// 获取服务的连接地址
	GetConn(string) (string, error)
}
