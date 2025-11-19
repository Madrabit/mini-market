package common

import (
	"github.com/kelseyhightower/envconfig"
	"os"
	"strings"
)

type Config struct {
	DB             DBConfig
	Server         ServerConfig
	LogLevel       string
	LogDevelopMode bool
	AllowedOrigins []string
	TestEmail      string
	TestMessage    string
}

type DBConfig struct {
	Server   string `envconfig:"SERVER" required:"true"`
	Port     int    `envconfig:"PORT" required:"true"`
	User     string `envconfig:"USER" required:"true"`
	Pass     string `envconfig:"PASS" required:"true"`
	Database string `envconfig:"DATABASE" required:"true"`
}

type ServerConfig struct {
	Address string `envconfig:"ADDRESS" required:"true"`
	Port    string `envconfig:"PORT" required:"true"`
}

func Load() (Config, error) {
	var cfg Config = Config{
		LogLevel:       os.Getenv("LOG_LEVEL"),
		LogDevelopMode: os.Getenv("LOG_DEVELOP_MODE") == "true",
		TestEmail:      os.Getenv("TEST_EMAIL"),
		TestMessage:    os.Getenv("TEST_MESSAGE"),
		AllowedOrigins: strings.Split(os.Getenv("CORS_ORIGINS"), ","),
	}
	if db, err := LoadDbConfig(); err != nil {
		return Config{}, err
	} else {
		cfg.DB = db
	}
	if server, err := LoadServerConfig(); err != nil {
		return Config{}, err
	} else {
		cfg.Server = server
	}
	return cfg, nil
}

func LoadDbConfig() (DBConfig, error) {
	var cfg DBConfig
	err := envconfig.Process("DB", &cfg)
	if err != nil {
		return DBConfig{}, err
	}
	return cfg, nil
}

func LoadServerConfig() (ServerConfig, error) {
	var cfg ServerConfig
	err := envconfig.Process("SERVER", &cfg)
	if err != nil {
		return ServerConfig{}, err
	}
	return cfg, nil
}
