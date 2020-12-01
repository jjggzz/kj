package consul

import (
	"errors"
	"github.com/jjggzz/kj/discovery"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"reflect"
	"time"
)

func (c *Client) DiscoveryService(servers []discovery.RpcServer) error {
	for _, serv := range servers {
		err := c.discoveryFromConsul(serv)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) GetConn(server discovery.RpcServer) (*grpc.ClientConn, error) {
	serviceName := reflect.TypeOf(server).String()
	instances := c.serverTable[serviceName]
	if len(instances) == 0 {
		return nil, errors.New("没有此服务")
	}
	rand.Seed(time.Now().Unix())
	node := rand.Int() % len(instances)
	log.Printf("连接到服务[%s]的第[%d]节点:[%s]", serviceName, node, instances[node].Address)
	return instances[node].Conn, nil
}

// 服务发现
func (c *Client) discoveryFromConsul(server discovery.RpcServer) error {
	serviceName := reflect.TypeOf(server).String()
	entity, _, err := c.consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return err
	}
	for _, value := range entity {
		conn, err := grpc.Dial(value.Service.ID, grpc.WithInsecure())
		if err != nil {
			log.Printf("连接[%s]服务的节点[%s]失败..", serviceName, value.Service.ID)
		}
		log.Printf("连接[%s]服务的节点[%s]成功..", serviceName, value.Service.ID)
		instance := discovery.Instance{Address: value.Service.ID, Conn: conn}
		c.serverTable[serviceName] = append(c.serverTable[serviceName], instance)
	}
	return nil
}
