package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/jjggzz/kj/discovery"
	"log"
	"reflect"
)

func (c *Client) RegisterService(servers []discovery.RpcServer) error {
	for _, server := range servers {
		err := c.registerToConsul(server)
		if err != nil {
			return err
		}
	}
	return nil
}

// 服务注册
func (c *Client) registerToConsul(server discovery.RpcServer) error {
	serviceName := reflect.TypeOf(server).String()
	// 配置注册到consul的服务的信息
	registration := new(api.AgentServiceRegistration)
	registration.ID = fmt.Sprintf("%s:%d", c.Config.Server.Ip, c.Config.Server.Tcp.Port)
	registration.Name = serviceName
	registration.Port = c.Config.Server.Tcp.Port
	registration.Address = c.Config.Server.Ip
	c.healthCheck(registration)
	log.Printf("开始注册服务[%s]到[%s]...", serviceName, c.Config.Discovery.Consul.Address)
	err := c.consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		log.Printf("服务[%s]注册失败", serviceName)
		panic(err)
	}
	log.Printf("服务[%s]注册成功", serviceName)
	return nil
}

// 健康检查
func (c *Client) healthCheck(registration *api.AgentServiceRegistration) {
	check := new(api.AgentServiceCheck)
	check.TCP = fmt.Sprintf("%s:%d", registration.Address, registration.Port)
	check.Timeout = fmt.Sprintf("%ds", c.Config.Discovery.Consul.Health.Timeout)
	check.Interval = fmt.Sprintf("%ds", c.Config.Discovery.Consul.Health.Interval)
	// 故障检查失败30s后 consul自动将注册服务删除
	check.DeregisterCriticalServiceAfter = fmt.Sprintf("%ds", c.Config.Discovery.Consul.Health.DeregisterCriticalServiceAfter)
	registration.Check = check
}
