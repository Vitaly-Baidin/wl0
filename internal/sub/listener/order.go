package listener

import (
	"context"
	"github.com/Vitaly-Baidin/l0/internal/sub/service"
	"github.com/Vitaly-Baidin/l0/logging"
	"github.com/Vitaly-Baidin/l0/pkg/util"
	"github.com/nats-io/stan.go"
)

type listener struct {
	logger       logging.Logger
	orderService *service.OrderService
}

func NewOrderListener(log logging.Logger, orderService *service.OrderService) *listener {
	return &listener{
		logger:       log,
		orderService: orderService,
	}
}

func (l *listener) StartListen(msg *stan.Msg) {
	order, err := util.ConvertJsonToOrder(msg.Data)
	if err != nil {
		l.logger.Errorf("failed convert message: %v", err)
		return
	}

	err = l.orderService.AddOrder(context.TODO(), order)
	if err != nil {
		l.logger.Errorf("failed save order: %v", err)
		return
	}
}
