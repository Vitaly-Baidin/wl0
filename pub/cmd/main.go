package main

import (
	"fmt"
	"github.com/Vitaly-Baidin/l0/pkg/config/viperconf"
	"github.com/Vitaly-Baidin/l0/pkg/logging/zaplog"
	"github.com/Vitaly-Baidin/l0/pkg/messageBroker/natsmb"
	"github.com/nats-io/stan.go"
	"os"
)

func main() {
	initLogger()

	cfg := initConfig()

	logger := zaplog.Logger

	model, err := os.ReadFile("model.json")
	if err != nil {
		logger.Errorf("failed to read file: %v\n", err)
		return
	}
	conn := connectStan(cfg)

	if err := conn.Publish("foo", model); err != nil {
		logger.Errorf("failed to publish message: %v\n", err)
	}

	err = conn.Close()
	if err != nil {
		logger.Errorf("failed to close publisher: %v\n", err)
	}
}

func initLogger() {
	fmt.Println("init logger")
	zaplog.InitializeLogger("logs/pub-service")
}

func initConfig() *viperconf.Config {
	logger := zaplog.Logger
	logger.Info("init config")
	loadConfig, err := viperconf.LoadConfig("pub/config")
	if err != nil {
		logger.Fatalf("failed to read config: %v\n", err)
	}
	logger.Info("init OK")
	return &loadConfig
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
