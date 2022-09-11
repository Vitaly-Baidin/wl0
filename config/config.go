package config

type Config struct {
	ServerConfig   ServerConfig   `yaml:"server"`
	DatabaseConfig DatabaseConfig `yaml:"database"`
	CacheConfig    CacheConfig    `yaml:"cache"`
	StanConfig     StanConfig     `yaml:"stan"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DatabaseConfig struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	DBName     string `yaml:"db_name"`
	Timeout    int    `yaml:"timeout"`
	MaxConnect int    `yaml:"max_connect"`
}

type CacheConfig struct {
	DefaultExpiration int `yaml:"default_expiration"`
	CleanupInterval   int `yaml:"cleanup_interval"`
}

type StanConfig struct {
	ClusterID string `yaml:"cluster_id"`
	ClientID  string `yaml:"client_id"`
	URL       string `yaml:"url"`
}
