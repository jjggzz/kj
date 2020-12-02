package consul

import (
	"errors"
	"fmt"
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

func (c *Client) GetConn(serverName string) (string, error) {
	// 在获取连接之前加上读锁，使其可以并发获取连接，而在异步更新服务列表时阻塞
	RWMutex.RLock()
	defer RWMutex.RUnlock()

	instances := c.serverTable[serverName]
	if len(instances) == 0 {
		return "", errors.New("此服务不存在可用节点")
	}
	rand.Seed(time.Now().Unix())
	node := rand.Int() % len(instances)
	log.Printf("访问服务[%s]的[%d]节点:[%s]", serverName, node, instances[node])
	return instances[node], nil
}

// 服务发现
func (c *Client) discoveryFromConsul(serverName string) error {
	entity, _, err := c.consulClient.Health().Service(serverName, "", true, nil)
	if err != nil {
		return err
	}
	for _, value := range entity {
		log.Printf("获取[%s]服务的节点成功:[%s]", serverName, value.Service.ID)
		c.serverTable[serverName] = append(c.serverTable[serverName], value.Service.ID)
	}
	return nil
}

// 更新服务列表
func (c *Client) updateServerList(serverNames []string) {
	for _, serverName := range serverNames {
		err := c.discoveryFromConsul(serverName)
		if err != nil {
			panic(fmt.Sprintf("异步更新[%s]服务节点列表失败", serverName))
		}
	}
}
