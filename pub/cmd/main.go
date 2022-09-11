package main

import (
	"fmt"
	"github.com/Vitaly-Baidin/l0/config"
	"github.com/Vitaly-Baidin/l0/config/viper"
	"github.com/Vitaly-Baidin/l0/logging"
	"github.com/Vitaly-Baidin/l0/logging/zapLog"
	"github.com/Vitaly-Baidin/l0/messageBroker/natsmb"
	"github.com/nats-io/stan.go"
	"log"
	"os"
)

func main() {
	logger := initLogger()

	cfg := initConfig(logger)

	model, err := os.ReadFile("model.json")
	if err != nil {
		logger.Fatalf("failed to read file: %v\n", err)
	}
	conn := connectStan(logger, cfg)

	if err := conn.Publish("order.message", model); err != nil {
		logger.Fatalf("failed to publish message: %v\n", err)
	}

	err = conn.Close()
	if err != nil {
		logger.Errorf("failed to close publisher: %v\n", err)
	}
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
	loadConfig, err := viper.LoadConfig("pub/config")
	if err != nil {
		log.Fatalf("error during configuration initialization: %v", err)
	}
	log.Info("initialization configuration successful")
	return *loadConfig
}

func connectStan(logger logging.Logger, cfg config.Config) stan.Conn {
	logger.Info("connect to nats-streaming")
	conn, err := natsmb.Connect(cfg)
	if err != nil {
		logger.Fatalf("failed to connect STAN aka nats-streaming: %v\n", err)
	}
	logger.Info("connect to nats-streaming OK")
	return conn
}
