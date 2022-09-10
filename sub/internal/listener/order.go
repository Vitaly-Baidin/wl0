package listener

import (
	"encoding/json"
	"github.com/Vitaly-Baidin/l0/pkg/logging/zaplog"
	"github.com/Vitaly-Baidin/l0/sub/internal/domain"
	"github.com/Vitaly-Baidin/l0/sub/internal/service"
	"github.com/nats-io/stan.go"
)

type listener struct {
	orderService *service.OrderService
	cacheService *service.CacheService
}

func NewOrderListener(orderService service.OrderService, cacheService service.CacheService) *listener {
	return &listener{
		orderService: &orderService,
		cacheService: &cacheService,
	}
}

func (l *listener) StartListen(msg *stan.Msg) {
	o := domain.Order{}

	err := json.Unmarshal(msg.Data, &o)
	if err != nil {
		zaplog.Logger.Errorf("invalid messege: %v\n", err)
		return
	}
	err = l.cacheService.SaveCache(*o.OrderUID, o)
	if err != nil {
		zaplog.Logger.Errorf("failed save to cache: %v\n", err)
		return
	}
	err = l.orderService.AddOrder(o)
	if err != nil {
		zaplog.Logger.Errorf("failed add order to db: %v\n", err)
		return
	}
}
