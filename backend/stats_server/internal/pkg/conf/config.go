package conf

import (
	"time"
	"net"
	"strconv"
)

const (
	ModeProd = "prod"
	ModeDev  = "dev"
)

type Config struct {
	Name string
	Host string
	Port int
	Mode string
}

type StatsConfig struct {
	CoinType               string
	MeasurementSharePrefix string
}

type PusherConfig struct {
	ReadInterval      time.Duration
	RetentionStrategy string
	Database          string
	Precision         string
	MeasurementPrefix string
	MeasurementSuffix string
	CoinType          string
	Mode              string
	PerSize           int
}

type GRPCConfig struct {
	Host string
	Port int
}

func (c *GRPCConfig) FormatHostPort() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}
