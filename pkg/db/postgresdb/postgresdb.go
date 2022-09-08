package postgresdb

import (
	"context"
	"fmt"
	"github.com/Vitaly-Baidin/l0/pkg/config/viperconf"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/url"
)

func NewPoolConfig(cfg *viperconf.Config) (*pgxpool.Config, error) {
	cfgDB := cfg.DatabaseConfig
	connStr :=
		fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
			"postgres",
			url.QueryEscape(cfgDB.Username),
			url.QueryEscape(cfgDB.Password),
			cfgDB.Host,
			cfgDB.Port,
			cfgDB.DbName,
			cfgDB.Timeout)
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	return poolConfig, nil
}

func NewConnection(poolConfig *pgxpool.Config, ctx context.Context) (*pgxpool.Pool, error) {
	conn, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
