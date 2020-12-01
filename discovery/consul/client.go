package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/jjggzz/kj/baseConfig"
	"github.com/jjggzz/kj/discovery"
	"log"
)

type Client struct {
	Config       baseConfig.Config
	consulClient *api.Client
	serverTable  map[string][]discovery.Instance
}

func NewConsulClient(conf baseConfig.Config) discovery.Discovery {
	// 创建consul客户端
	c := api.DefaultConfig()
	// 设置consul的地址
	c.Address = conf.Discovery.Consul.Address
	client, err := api.NewClient(c)
	if err != nil {
		log.Printf("连接consul失败: %s", err)
		panic(err)
	}
	return &Client{Config: conf, consulClient: client, serverTable: map[string][]discovery.Instance{}}
}
