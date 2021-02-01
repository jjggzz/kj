// author: JGZ
// time:   2021-01-28 18:06
package nacos

import (
	"github.com/jjggzz/kit/log"
	"github.com/jjggzz/kit/sd"
	sdnacos "github.com/jjggzz/kit/sd/nacos"
	"github.com/jjggzz/kj/discovery"
	"github.com/jjggzz/kj/uitls"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"strconv"
	"strings"
)

type nacos struct {
	address    string
	serverName string
	port       int
	namespace  string
	weight     float64
	logger     log.Logger
	client     naming_client.INamingClient
	r          sd.Registrar
}

func NewNacosDiscovery(address string, serverName string, port int, namespace string, weight float64, logger log.Logger) discovery.Discover {
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithRotateTime("1h"),
		constant.WithMaxAge(3),
		constant.WithLogLevel("debug"),
	)
	split := strings.Split(address, ":")
	p, err := strconv.Atoi(split[1])
	if err != nil {
		panic(err)
	}
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			split[0],
			uint64(p),
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		),
	}
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	n := nacos{
		client:     namingClient,
		address:    address,
		serverName: serverName,
		port:       port,
		namespace:  namespace,
		weight:     weight,
		logger:     logger,
	}

	return &n
}

func (n *nacos) RegisterServer() {
	param := sdnacos.Param{
		Ip:          uitls.LocalIpv4(),
		Port:        uint64(n.port),
		ServiceName: n.serverName,
		Weight:      n.weight,
		Metadata:    nil,
		ClusterName: "",
		GroupName:   "",
	}
	n.r = sdnacos.NewRegistrar(n.client, param, n.logger)
	n.r.Register()
}

func (n *nacos) DeregisterServer() {
	n.r.Deregister()
}

func (n *nacos) Discovery(targetServerName string) (sd.Instancer, error) {
	instancer := sdnacos.NewInstancer(n.client, n.logger, targetServerName, "", nil)
	return instancer, nil
}
