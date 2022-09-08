package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Vitaly-Baidin/l0/pkg/config/viperconf"
	"github.com/Vitaly-Baidin/l0/pkg/db/postgresdb"
	"github.com/Vitaly-Baidin/l0/pkg/logging/zaplog"
	"github.com/Vitaly-Baidin/l0/pkg/messageBroker/natsmb"
	myCache "github.com/Vitaly-Baidin/l0/sub/internal/cache"
	_ "github.com/Vitaly-Baidin/l0/sub/internal/migrations"
	"github.com/Vitaly-Baidin/l0/sub/internal/order"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
	"github.com/pressly/goose/v3"
	"os"
	"os/signal"
	"time"
)

func main() {
	initLogger()

	cfg := initConfig()

	poolConfig := createPoolConfig(cfg)
	connDB := connectDB(poolConfig)
	defer connDB.Close()
	initMigrations(poolConfig)

	connStan := connectStan(cfg)
	defer connStan.Close()

	cacheService := initCache(cfg, connDB)

	orderService := order.NewOrderService(connDB)

	listener := order.NewOrderListener(*orderService, *cacheService)

	sub, err := connStan.Subscribe("foo", listener.StartListen) // TODO вынести в отдельный сервис

	defer sub.Unsubscribe()
	if err != nil {
		connStan.Close()
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing ALL connection...\n\n")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func initLogger() {
	fmt.Println("init logger")
	zaplog.InitializeLogger("logs/sub-service")
	fmt.Println("init logger OK")
}

func initConfig() *viperconf.Config {
	logger := zaplog.Logger
	logger.Info("init config")
	loadConfig, err := viperconf.LoadConfig("sub/config")
	if err != nil {
		logger.Fatalf("failed to read config %v", err)
	}
	logger.Info("init config OK")
	return &loadConfig
}

func createPoolConfig(cfg *viperconf.Config) *pgxpool.Config {
	logger := zaplog.Logger
	logger.Info("create pool config")
	poolConfig, err := postgresdb.NewPoolConfig(cfg)
	if err != nil {
		logger.Fatalf("failed to create pool config %v\n", err)
	}

	logger.Info("create pool config OK")
	poolConfig.MaxConns = 5
	return poolConfig
}

func connectDB(poolConfig *pgxpool.Config) *pgxpool.Pool {
	logger := zaplog.Logger
	logger.Info("connect to database")

	c, err := postgresdb.NewConnection(poolConfig)
	if err != nil {
		logger.Fatalf("failed connect to database: %v\n", err)
	}

	logger.Info("connect to database OK")

	logger.Info("Ping database")
	_, err = c.Exec(context.Background(), ";")
	if err != nil {
		logger.Errorf("Ping failed: %v\n", err)
	}
	logger.Info("Ping database OK")

	return c
}

func initCache(cfg *viperconf.Config, database *pgxpool.Pool) *myCache.Service {
	logger := zaplog.Logger
	logger.Info("init cache")
	DefaultExpiration := time.Duration(cfg.CacheConfig.DefaultExpiration) * time.Minute
	CleanupInterval := time.Duration(cfg.CacheConfig.CleanupInterval) * time.Minute

	c := cache.New(DefaultExpiration, CleanupInterval)
	service := myCache.NewCacheService(database, c)

	c.OnEvicted(func(key string, v interface{}) {
		service.RemoveCacheFromDB(key)
	})

	logger.Info("init cache OK")

	logger.Info("load cache from db")
	cashes := service.GetAllCacheFromDB()

	for _, cash := range cashes {
		c.Set(cash.Key, cash.Value, cash.Expiration)
	}
	logger.Info("load cache from db OK")

	return service
}

func initMigrations(poolConfig *pgxpool.Config) {
	logger := zaplog.Logger
	logger.Info("migrate database")

	mdb, _ := sql.Open("postgres", poolConfig.ConnString())
	err := goose.Up(mdb, "/var")
	if err != nil {
		zaplog.Logger.Errorf("failed to migrate db: %v\n", err)
	}

	logger.Info("migrate database OK")
}

func connectStan(cfg *viperconf.Config) stan.Conn {
	logger := zaplog.Logger
	logger.Info("connect to nats-streaming")
	conn, err := natsmb.Connect(cfg)
	if err != nil {
		logger.Fatalf("failed to connect STAN aka nats-streaming: %v\n", err)
	}
	logger.Info("connect to nats-streaming OK")
	return conn
}
