package discovery

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	sdconsul "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/jjggzz/kj/uitls"
	"os"
)

type consul struct {
	address    string
	serverName string
	port       int
	logger     log.Logger
	r          sd.Registrar
}

func NewConsulDiscovery(address string, serverName string, port int, logger log.Logger) Discover {
	return &consul{address: address, serverName: serverName, port: port, logger: logger}
}

func (c *consul) RegisterServer() {
	var client sdconsul.Client
	{
		consulCfg := api.DefaultConfig()
		consulCfg.Address = c.address
		consulClient, err := api.NewClient(consulCfg)
		if err != nil {
			_ = c.logger.Log("create consul client error:", err)
			os.Exit(1)
		}
		client = sdconsul.NewClient(consulClient)
	}

	// 健康检测
	check := api.AgentServiceCheck{
		TCP:                            fmt.Sprintf("%s:%d", uitls.LocalIpv4(), c.port),
		Interval:                       "10s",
		Timeout:                        "1s",
		Notes:                          "Consul check service health status.",
		DeregisterCriticalServiceAfter: "30s",
	}
	// 服务名
	reg := api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s:%d", uitls.LocalIpv4(), c.port),
		Name:    c.serverName,
		Address: uitls.LocalIpv4(),
		Port:    c.port,
		Check:   &check,
	}
	c.r = sdconsul.NewRegistrar(client, &reg, c.logger)
	c.r.Register()
}

func (c *consul) DeregisterServer() {
	c.r.Deregister()
}
