package natsmb

import (
	"github.com/Vitaly-Baidin/l0/pkg/config/viperconf"
	"github.com/nats-io/stan.go"
)

func Connect(cfg *viperconf.Config) (stan.Conn, error) {
	return stan.Connect(
		cfg.StanConfig.ClusterID,
		cfg.StanConfig.ClientID,
		stan.NatsURL(cfg.StanConfig.URL),
	)
}
