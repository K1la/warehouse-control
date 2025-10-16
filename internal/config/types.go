package config

import "time"

type Config struct {
	Postgres   Postgres   `mapstructure:"postgres"`
	HTTPServer HTTPServer `mapstructure:"http_server"`
	JWT        JWT        `mapstructure:"jwt"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
}

type HTTPServer struct {
	Address     string `mapstructure:"address"`
	Timeout     int    `mapstructure:"timeout"`
	IdleTimeout int    `mapstructure:"idle_timeout"`
}

type JWT struct {
	Secret string        `mapstructure:"secret"`
	TTL    time.Duration `mapstructure:"ttl"`
}
