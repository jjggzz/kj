package consul

import (
	"errors"
	"fmt"
	"github.com/jjggzz/kj/discovery"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"time"
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
				c.updateServerList(serverNames)
			}
		}
	}()

	return nil
}

func (c *Client) GetOneConn(serverName string) (*grpc.ClientConn, error) {
	// 在获取连接之前加上读锁，使其可以并发获取连接，而在异步更新服务列表时阻塞
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	instances := c.serverTable[serverName]
	if len(instances) == 0 {
		return nil, errors.New("此服务不存在可用节点")
	}
	rand.Seed(time.Now().Unix())
	node := rand.Int() % len(instances)
	log.Printf("获取服务[%s]的[%d]节点:[%s]", serverName, node, instances[node].Address)
	return instances[node].Conn, nil
}

// 根据服务名获取服务连接列表
func (c *Client) GetConn(serverName string) ([]*grpc.ClientConn, error) {
	// 在获取连接之前加上读锁，使其可以并发获取连接，而在异步更新服务列表时阻塞
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	conns := make([]*grpc.ClientConn, len(c.serverTable[serverName]))
	for i, v := range c.serverTable[serverName] {
		conns[i] = v.Conn
	}
	if len(conns) == 0 {
		return nil, errors.New("暂无可用连接")
	}
	return conns, nil
}

// 服务发现
func (c *Client) discoveryFromConsul(serverName string) error {
	entity, _, err := c.consulClient.Health().Service(serverName, "", true, nil)
	if err != nil {
		return err
	}
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	for _, v := range c.serverTable[serverName] {
		_ = v.Conn.Close()
	}
	c.serverTable[serverName] = make([]discovery.Instance, 0)
	for _, value := range entity {
		conn, err := grpc.Dial(value.Service.ID, grpc.WithInsecure())
		if err != nil {
			log.Printf("连接[%s]服务的节点失败:[%s]", serverName, value.Service.ID)
			continue
		}
		//log.Printf("连接[%s]服务的节点成功:[%s]", serverName, value.Service.ID)
		c.serverTable[serverName] = append(c.serverTable[serverName], discovery.Instance{Address: value.Service.ID, Conn: conn})
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
