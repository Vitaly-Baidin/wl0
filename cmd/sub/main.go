package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Vitaly-Baidin/l0/config"
	"github.com/Vitaly-Baidin/l0/config/viper"
	"github.com/Vitaly-Baidin/l0/db/postgresdb"
	"github.com/Vitaly-Baidin/l0/internal/sub/listener"
	_ "github.com/Vitaly-Baidin/l0/internal/sub/migrations"
	"github.com/Vitaly-Baidin/l0/internal/sub/route"
	"github.com/Vitaly-Baidin/l0/internal/sub/service"
	"github.com/Vitaly-Baidin/l0/logging"
	"github.com/Vitaly-Baidin/l0/logging/zapLog"
	"github.com/Vitaly-Baidin/l0/messageBroker/natsmb"
	"github.com/Vitaly-Baidin/l0/pkg/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	logger := initLogger()
	cfg := initConfig(logger)

	ctx := context.Background()
	connDB := connectDB(ctx, logger, cfg)
	initMigrations(logger, connDB)

	connStan := connectStan(logger, cfg)

	cacheService := initCacheService(ctx, logger, cfg, connDB)

	orderService := initOrderCache(cacheService, connDB)

	orderListener := listener.NewOrderListener(logger, orderService)

	_, err := connStan.Subscribe("order.message", orderListener.StartListen, stan.DurableName("last"))
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
		logger.Info("connect to server")
		err = http.ListenAndServe(url, r)
		if err != nil {
			logger.Fatalf("failed to create server: %v\n", err)
			return
		}
	}()
	// close operation
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			logger.Info("Received an interrupt, closing ALL connection...")
			err = connStan.Close()
			if err != nil {
				logger.Errorf("failed to close stan connect %v\n", err)
			}
			connDB.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func initLogger() logging.Logger {
	fmt.Println("initialization logger...")
	logger, err := zapLog.InitializeLogger("logs/sub-service")
	if err != nil {
		log.Fatalf("error during logger initialization: %v", err)
	}
	fmt.Println("initialization logger successful")
	return logger
}

func initConfig(log logging.Logger) config.Config {
	log.Info("initialization configuration...")
	loadConfig, err := viper.LoadConfig(".")
	if err != nil {
		log.Fatalf("error during configuration initialization: %v", err)
	}
	log.Info("initialization configuration successful")
	return *loadConfig
}

func connectDB(ctx context.Context, log logging.Logger, cfg config.Config) *pgxpool.Pool {
	log.Info("connection to database...")
	c, err := postgresdb.NewConnection(ctx, cfg)
	if err != nil {
		log.Fatalf("failed connect to database: %v\n", err)
	}
	log.Info("connection to database successful")

	return c
}

func initMigrations(log logging.Logger, poolConfig *pgxpool.Pool) {
	log.Info("check migration database...")

	mdb, _ := sql.Open("postgres", poolConfig.Config().ConnString())
	err := goose.Up(mdb, "/var")
	if err != nil {
		log.Errorf("failed to migrate db: %v\n", err)
	}

	log.Info("check migration database successful")
}

func connectStan(log logging.Logger, cfg config.Config) stan.Conn {
	log.Info("connection to nats-streaming...")
	conn, err := natsmb.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect STAN aka nats-streaming: %v\n", err)
	}
	log.Info("connection to nats-streaming successful")
	return conn
}

func initCacheService(
	ctx context.Context,
	log logging.Logger,
	cfg config.Config,
	database *pgxpool.Pool) *service.CacheService {
	log.Info("initialization cache...")

	var cacheRepository repository.CacheRepository = repository.NewCacheRepository(database)
	cacheService := service.NewCacheService(cacheRepository)
	cacheService.InitCache(log, cfg)

	log.Info("initialization cache successful")

	log.Info("load cache from db")
	err := cacheService.LoadCacheFromDB(ctx)
	if err != nil {
		log.Infof("failed load cache from db: %v", err)
		return nil
	}

	log.Info("load cache from db successful")

	return cacheService
}

func initOrderCache(cacheService *service.CacheService, database *pgxpool.Pool) *service.OrderService {
	var orderRepository repository.OrderRepository = repository.NewOrderRepository(database)

	return service.NewOrderService(orderRepository, cacheService)
}
