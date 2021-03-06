package conf

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	ModeProd = "prod"
	ModeDev  = "dev"
)

type BTCClientConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	UseSSL   bool
}

type RedisConfig struct {
	Host  string
	Port  int
	Password string
	DBNum int
	RedisPoolConfig
}

type RedisPoolConfig struct {
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

type CoinConfig struct {
	PullGBTInterval  time.Duration
	NotifyInterval   time.Duration
	JobCacheExpireTs time.Duration
}

type SphereConfig struct {
	Configs map[string]CoinConfig
}

type ServerConfig struct {
	Name string
	Host string
	Port int
	Mode string
}

func (c *ServerConfig) IsValidMode() bool {
	if c.Mode == ModeProd {
		return true
	}

	if c.Mode == ModeDev {
		return true
	}
	return false
}

func NewSphereConfig() *SphereConfig {
	configs := make(map[string]CoinConfig, 0)
	return &SphereConfig{Configs:configs}
}


func (c *ServerConfig) FormatHostPort() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}


func (conf *RedisConfig) FormatAddress() string {
	return fmt.Sprintf("%s:%d", conf.Host, conf.Port)
}

