package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/jjggzz/kj/uitls"
	"log"
)

// 服务注册
func (c *Client) RegisterServer() error {
	ipv4 := uitls.LocalIpv4()
	// 配置注册到consul的服务的信息
	registration := new(api.AgentServiceRegistration)
	registration.ID = fmt.Sprintf("%s:%d", ipv4, c.ser.Tcp.Port)
	registration.Name = c.ser.ServerName
	registration.Port = c.ser.Tcp.Port
	registration.Address = ipv4
	c.healthCheck(registration)
	err := c.consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		log.Printf("服务[%s]注册到[%s]失败", c.ser.ServerName, c.dis.Consul.Address)
		return err
	}
	log.Printf("服务[%s]注册到[%s]成功", c.ser.ServerName, c.dis.Consul.Address)
	return nil
}

// 健康检查
func (c *Client) healthCheck(registration *api.AgentServiceRegistration) {
	check := new(api.AgentServiceCheck)
	check.TCP = fmt.Sprintf("%s:%d", registration.Address, registration.Port)
	check.Timeout = fmt.Sprintf("%ds", c.dis.Consul.Health.Timeout)
	check.Interval = fmt.Sprintf("%ds", c.dis.Consul.Health.Interval)
	// 故障检查失败30s后 consul自动将注册服务删除
	check.DeregisterCriticalServiceAfter = fmt.Sprintf("%ds", c.dis.Consul.Health.DeregisterCriticalServiceAfter)
	registration.Check = check
}
