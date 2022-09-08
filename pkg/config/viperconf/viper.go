package viperconf

import (
	"github.com/Vitaly-Baidin/l0/pkg/logging/zaplog"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig   ServerConfig   `mapstructure:"server"`
	DatabaseConfig DatabaseConfig `mapstructure:"database"`
	CacheConfig    CacheConfig    `mapstructure:"cache"`
	StanConfig     StanConfig     `mapstructure:"stan"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DbName   string `mapstructure:"db_name"`
	Timeout  int    `mapstructure:"timeout"`
}

type CacheConfig struct {
	DefaultExpiration int32 `mapstructure:"default_expiration"`
	CleanupInterval   int32 `mapstructure:"cleanup_interval"`
}

type StanConfig struct {
	ClusterID string `mapstructure:"cluster_id"`
	ClientID  string `mapstructure:"client_id"`
	URL       string `mapstructure:"url"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		zaplog.Logger.Panic("failed to read config")
		return
	}

	serverDefault()
	cacheDefault()
	stanDefault()

	err = viper.Unmarshal(&config)
	return
}

func serverDefault() {
	viper.SetDefault("server.port", "8080")
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
