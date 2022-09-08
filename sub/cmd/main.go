package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Vitaly-Baidin/l0/pkg/config/viperconf"
	"github.com/Vitaly-Baidin/l0/pkg/db/postgresdb"
	"github.com/Vitaly-Baidin/l0/pkg/logging/zaplog"
	"github.com/Vitaly-Baidin/l0/pkg/messageBroker/natsmb"
	"github.com/Vitaly-Baidin/l0/sub/internal/listener"
	_ "github.com/Vitaly-Baidin/l0/sub/internal/migrations"
	"github.com/Vitaly-Baidin/l0/sub/internal/route"
	"github.com/Vitaly-Baidin/l0/sub/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/pressly/goose/v3"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	initLogger()

	cfg := initConfig()

	ctx := context.Background()
	ctxValue := context.WithValue(ctx, "cache.expiration", cfg.CacheConfig.DefaultExpiration)
	ctxTimeOut, cancel := context.WithCancel(ctxValue)
	defer cancel()

	poolConfig := createPoolConfig(cfg)
	connDB := connectDB(poolConfig, ctxTimeOut)
	initMigrations(poolConfig)

	connStan := connectStan(cfg)

	cacheService := initCache(cfg, ctxTimeOut, connDB)

	orderService := service.NewOrderService(connDB, ctxTimeOut, cacheService)

	listener := listener.NewOrderListener(*orderService, *cacheService)

	_, err := connStan.Subscribe("order.message", listener.StartListen, stan.DurableName("last"))
	if err != nil {
		connStan.Close()
		connDB.Close()
		return
	}

	r := chi.NewRouter()

	orderRoute := route.NewOrderRoute(orderService)
	orderRoute.RegisterRoute(r)

	go func() {
		url := fmt.Sprintf("%s:%d", cfg.ServerConfig.Host, cfg.ServerConfig.Port)
		zaplog.Logger.Info("connect to server")
		err = http.ListenAndServe(url, r)
		if err != nil {
			zaplog.Logger.Fatalf("failed to create server: %v\n", err)
			return
		}
	}()
	// close operation
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			zaplog.Logger.Info("Received an interrupt, closing ALL connection...")
			err := connStan.Close()
			if err != nil {
				zaplog.Logger.Errorf("failed to close stan connect %v\n", err)
			}
			connDB.Close()
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
	poolConfig.MaxConns = 10
	return poolConfig
}

func connectDB(poolConfig *pgxpool.Config, ctx context.Context) *pgxpool.Pool {
	logger := zaplog.Logger
	logger.Info("connect to database")

	c, err := postgresdb.NewConnection(poolConfig, ctx)
	if err != nil {
		logger.Fatalf("failed connect to database: %v\n", err)
	}

	logger.Info("connect to database OK")

	logger.Info("Ping database")
	for i := 0; i < 10; i++ {
		go func() {
			_, err = c.Exec(context.Background(), ";")
		}()
	}
	if err != nil {
		logger.Errorf("Ping failed: %v\n", err)
	}
	logger.Info("Ping database OK")

	return c
}

func initCache(cfg *viperconf.Config, ctx context.Context, database *pgxpool.Pool) *service.CacheService {
	logger := zaplog.Logger
	logger.Info("init cache")

	service := service.NewCacheService(database, ctx, cfg)
	service.Cache.OnEvicted(func(key string, v interface{}) {
		err := service.RemoveCacheFromDB(key)
		if err != nil {
			logger.Errorf("failed remove item cache from db: %v\n", err)
			return
		}
	})

	logger.Info("init cache OK")

	logger.Info("load cache from db")
	caches, err := service.GetAllCacheFromDB()
	if errors.Is(err, pgx.ErrNoRows) {
		logger.Warn("cache in db is empty")
		logger.Info("load cache from db OK")
		return service
	} else if err != nil {
		logger.Fatalf("failed to load cache from db: %v\n", err)
	}

	for _, v := range caches {
		service.Cache.Set(v.Key, v.Value, v.Expiration*time.Minute)
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
