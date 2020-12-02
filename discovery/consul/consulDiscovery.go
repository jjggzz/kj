package consul

import (
	"errors"
	"github.com/jjggzz/kj/discovery"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"reflect"
	"sync"
	"time"
)

var (
	RWMutex sync.RWMutex
)

func (c *Client) DiscoveryService(servers []discovery.RpcServer) error {
	for _, serv := range servers {
		err := c.discoveryFromConsul(serv)
		if err != nil {
			return err
		}
	}

	// 异步函数执行服务刷新
	go func() {
		for {
			time.Sleep(time.Second * 5)
			RWMutex.Lock()
			c.updateServerList(servers)
			RWMutex.Unlock()
		}
	}()

	return nil
}

func (c *Client) GetConn(server discovery.RpcServer) (*grpc.ClientConn, error) {
	serviceName := reflect.TypeOf(server).String()
	// 在获取连接之前加上读锁，使其可以并发获取连接，而在异步更新服务列表时阻塞
	RWMutex.RLock()
	defer RWMutex.RUnlock()

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

// 更新服务列表
func (c *Client) updateServerList(servers []discovery.RpcServer) {
	for _, server := range servers {
		serverName := reflect.TypeOf(server).String()
		entity, _, err := c.consulClient.Health().Service(serverName, "", true, nil)
		if err != nil {
			panic("异步更新服务列表失败")
		}
		c.serverTable = make(map[string][]discovery.Instance)
		for _, value := range entity {
			conn, err := grpc.Dial(value.Service.ID, grpc.WithInsecure())
			if err != nil {
				log.Printf("连接[%s]服务的节点[%s]失败..", serverName, value.Service.ID)
			}
			instance := discovery.Instance{Address: value.Service.ID, Conn: conn}
			c.serverTable[serverName] = append(c.serverTable[serverName], instance)
		}

	}
}
