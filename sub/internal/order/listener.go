package order

import (
	"encoding/json"
	"github.com/Vitaly-Baidin/l0/pkg/logging/zaplog"
	myCache "github.com/Vitaly-Baidin/l0/sub/internal/cache"
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
)

type listener struct {
	orderService *Service
	cacheService *myCache.Service
}

func NewOrderListener(orderService Service, cacheService myCache.Service) *listener {
	return &listener{
		orderService: &orderService,
		cacheService: &cacheService,
	}
}

func (l *listener) StartListen(msg *stan.Msg) {
	o := Order{}

	err := json.Unmarshal(msg.Data, &o)
	if err != nil {
		zaplog.Logger.Errorf("invalid messege: %v\n", err)
	}
	l.cacheService.Cache.Set("foo", o, cache.DefaultExpiration)
	l.cacheService.SaveCache(o.OrderUID, o, cache.DefaultExpiration)
	l.orderService.AddOrder(o)
}
