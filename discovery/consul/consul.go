package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/jjggzz/kit/log"
	"github.com/jjggzz/kit/sd"
	sdconsul "github.com/jjggzz/kit/sd/consul"
	"github.com/jjggzz/kj/discovery"
	"github.com/jjggzz/kj/uitls"
	"os"
)

type consul struct {
	address    string
	serverName string
	port       int
	logger     log.Logger
	client     sdconsul.Client
	r          sd.Registrar
}

func NewConsulDiscovery(address string, serverName string, port int, logger log.Logger) discovery.Discover {
	consulCfg := api.DefaultConfig()
	consulCfg.Address = address
	consulClient, err := api.NewClient(consulCfg)
	if err != nil {
		_ = logger.Log("create consul client error:", err)
		os.Exit(1)
	}
	client := sdconsul.NewClient(consulClient)
	return &consul{address: address, serverName: serverName, port: port, logger: logger, client: client}
}

func (c *consul) RegisterServer() {
	// 健康检测
	check := api.AgentServiceCheck{
		TCP:                            fmt.Sprintf("%s:%d", uitls.LocalIpv4(), c.port),
		Interval:                       "5s",
		Timeout:                        "5s",
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
	c.r = sdconsul.NewRegistrar(c.client, &reg, c.logger)
	c.r.Register()
}

func (c *consul) DeregisterServer() {
	c.r.Deregister()
}

func (c *consul) Discovery(targetServerName string) (sd.Instancer, error) {
	instancer := sdconsul.NewInstancer(c.client, c.logger, targetServerName, []string{}, true)
	return instancer, nil
}
