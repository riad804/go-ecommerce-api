package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Redis   RedisConfig   `mapstructure:"redis"`
	Token   TokenConfig   `mapstructure:"token"`
	MongoDB MongoDBConfig `mapstructure:"mongodb"`
	Email   EmailConfig   `mapstructure:"email"`
}

type EmailConfig struct {
	Name       string        `mapstructure:"sender_name"`
	Address    string        `mapstructure:"sender_address"`
	Password   string        `mapstructure:"sender_password"`
	Expiration time.Duration `mapstructure:"exp_duration"`
}

type MongoDBConfig struct {
	URI      string        `mapstructure:"uri"`
	Database string        `mapstructure:"database"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type TokenConfig struct {
	AccessKey   string        `mapstructure:"access_key"`
	AccDuration time.Duration `mapstructure:"access_duration"`
	RefreshKey  string        `mapstructure:"refresh_key"`
	RefDuration time.Duration `mapstructure:"refresh_duration"`
}

type ServerConfig struct {
	Name      string          `mapstructure:"name"`
	Port      int             `mapstructure:"port"`
	Host      string          `mapstructure:"host"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

type RateLimitConfig struct {
	Enabled         bool          `mapstructure:"enabled"`
	RateLimit       int           `mapstructure:"rate_limit"`
	RateLimitWindow time.Duration `mapstructure:"rate_limit_window"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
