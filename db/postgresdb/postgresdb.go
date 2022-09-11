package postgresdb

import (
	"context"
	"fmt"
	"github.com/Vitaly-Baidin/l0/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/url"
)

func NewConnection(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	cfgDB := cfg.DatabaseConfig
	connStr :=
		fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
			"postgres",
			url.QueryEscape(cfgDB.Username),
			url.QueryEscape(cfgDB.Password),
			cfgDB.Host,
			cfgDB.Port,
			cfgDB.DBName,
			cfgDB.Timeout)
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("error parse config file: %v", err)
	}
	poolConfig.MaxConns = int32(cfg.DatabaseConfig.MaxConnect)

	conn, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error connect to db: %v", err)
	}
	return conn, nil
}
