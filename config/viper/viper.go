package viper

import (
	"fmt"
	"github.com/Vitaly-Baidin/l0/config"
	"github.com/mitchellh/mapstructure"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
)

func LoadConfig(path string) (config *config.Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	serverDefault()
	databaseDefault()
	cacheDefault()
	stanDefault()

	err = viper.Unmarshal(&config, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "yaml"
	})
	return config, nil
}

func serverDefault() {
	viper.SetDefault("server.port", 8080)
}

func databaseDefault() {
	viper.SetDefault("database.timeout", 5)
	viper.SetDefault("database.max_connect", 1)
}

func cacheDefault() {
	viper.SetDefault("cache.default_expiration", 5)
	viper.SetDefault("cache.cleanup_interval", 10)
}

func stanDefault() {
	viper.SetDefault("stan.cluster_id", "test-cluster")
	viper.SetDefault("stan.client_id", "test-client")
	viper.SetDefault("stan.url", stan.DefaultNatsURL)
}
