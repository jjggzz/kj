package consul

import (
	"errors"
	"github.com/jjggzz/kj/discovery"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	RWMutex sync.RWMutex
)

func (c *Client) DiscoveryServers(serverNames []string) error {
	for _, serverName := range serverNames {
		err := c.discoveryFromConsul(serverName)
		if err != nil {
			return err
		}
	}

	// 异步函数执行服务刷新
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-ticker.C:
				RWMutex.Lock()
				c.updateServerList(serverNames)
				RWMutex.Unlock()
			}
		}
	}()

	return nil
}

func (c *Client) GetConn(serverName string) (*grpc.ClientConn, error) {
	// 在获取连接之前加上读锁，使其可以并发获取连接，而在异步更新服务列表时阻塞
	RWMutex.RLock()
	defer RWMutex.RUnlock()

	instances := c.serverTable[serverName]
	if len(instances) == 0 {
		return nil, errors.New("没有此服务")
	}
	rand.Seed(time.Now().Unix())
	node := rand.Int() % len(instances)
	log.Printf("访问服务[%s]的第[%d]节点:[%s]", serverName, node, instances[node].Address)
	return instances[node].Conn, nil
}

// 服务发现
func (c *Client) discoveryFromConsul(serverName string) error {
	entity, _, err := c.consulClient.Health().Service(serverName, "", true, nil)
	if err != nil {
		return err
	}
	for _, value := range entity {
		conn, err := grpc.Dial(value.Service.ID, grpc.WithInsecure())
		if err != nil {
			log.Printf("连接[%s]服务的节点[%s]失败..", serverName, value.Service.ID)
		}
		log.Printf("连接[%s]服务的节点[%s]成功..", serverName, value.Service.ID)
		instance := discovery.Instance{Address: value.Service.ID, Conn: conn}
		c.serverTable[serverName] = append(c.serverTable[serverName], instance)
	}
	return nil
}

// 更新服务列表
func (c *Client) updateServerList(serverNames []string) {
	for _, serverName := range serverNames {
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
