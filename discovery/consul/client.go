package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/jjggzz/kj/config"
	"github.com/jjggzz/kj/discovery"
	"log"
	"sync"
)

type Client struct {
	ser          *config.Server
	dis          *config.Discovery
	consulClient *api.Client
	rwMutex      sync.RWMutex
	serverTable  map[string][]discovery.Instance
}

func NewConsulClient(ser *config.Server, dis *config.Discovery) discovery.Discover {
	// 创建consul客户端
	c := api.DefaultConfig()
	// 设置consul的地址
	c.Address = dis.Consul.Address
	client, err := api.NewClient(c)
	if err != nil {
		log.Printf("连接consul失败: %s", err)
		panic(err)
	}
	return &Client{ser: ser, dis: dis, consulClient: client, serverTable: map[string][]discovery.Instance{}}
}
