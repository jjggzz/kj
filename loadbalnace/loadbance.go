package loadbalnace

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"time"
)

const (
	Round  = 1
	Random = 2
)

// strategy:负载策略 1轮询(默认) 2随机
func BuildLoadBalance(instancer sd.Instancer, factory sd.Factory, strategy int, max int, timeout time.Duration, logger log.Logger) endpoint.Endpoint {
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	var balancer lb.Balancer
	switch strategy {
	case Random:
		balancer = lb.NewRandom(endpointer, time.Now().Unix())
	case Round:
		balancer = lb.NewRoundRobin(endpointer)
	default:
		balancer = lb.NewRoundRobin(endpointer)
	}
	retry := lb.Retry(max, timeout, balancer)
	return retry
}

func BuildDefaultLoadBalance(instancer sd.Instancer, factory sd.Factory, logger log.Logger) endpoint.Endpoint {
	return BuildLoadBalance(instancer, factory, Random, 3, 500*time.Millisecond, logger)
}
