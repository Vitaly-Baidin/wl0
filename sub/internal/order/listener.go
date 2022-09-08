package order

import (
	"context"
	"encoding/json"
	"github.com/Vitaly-Baidin/l0/pkg/logging/zaplog"
	"github.com/nats-io/stan.go"
)

type Listener struct {
	Repository *Repository
}

func (l *Listener) StartListen(msg *stan.Msg) {
	o := Order{}

	err := json.Unmarshal(msg.Data, &o)
	if err != nil {
		zaplog.Logger.Errorf("invalid messege: %v\n", err)
	}
	l.Repository.AddOrder(context.Background(), o)
}
