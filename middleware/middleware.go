package middleware

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"time"
)

// 日志中间件
func DefaultLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				_ = logger.Log("transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// 断路器中间件(该断路器只会返回错误)
// 配置为空则返回具有默认配置的断路器
// 在断路器半开的时候同事最多允许1个请求
// 断路器在半开状态下不清除计数
// 断路器半开时间默认为60s
// 5次失败就会处于断路状态
func BreakerMiddleware(setting gobreaker.Settings) endpoint.Middleware {
	return circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(setting))
}

const (
	LimitError = 1
	LimitDelay = 2
)

// 限流中间件
// count:每秒允许通过的请求量
// strategy:限流策略 1返回错误(默认) 2延迟通过
func LimitMiddleware(strategy int, count int) endpoint.Middleware {
	switch strategy {
	case LimitError:
		return ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), count))
	case LimitDelay:
		return ratelimit.NewDelayingLimiter(rate.NewLimiter(rate.Every(time.Second), count))
	default:
		return ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), count))
	}
}
