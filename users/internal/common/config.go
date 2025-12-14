package common

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
	"strings"
)

type Config struct {
	DB             DBConfig
	Server         ServerConfig
	LogLevel       string
	LogDevelopMode bool
	AppName        string
	AppVersion     string
	AllowedOrigins []string
}

type DBConfig struct {
	Host     string `envconfig:"HOST" required:"true"`
	Port     int    `envconfig:"PORT" required:"true"`
	User     string `envconfig:"USER" required:"true"`
	Pass     string `envconfig:"PASS" required:"true"`
	Name     string `envconfig:"NAME" required:"true"`
	Database string `envconfig:"DATABASE" required:"true"`
}
type ServerConfig struct {
	Address string `envconfig:"ADDRESS" required:"true"`
	Port    string `envconfig:"PORT" required:"true"`
}

func Load() (*Config, error) {
	var cfg Config = Config{
		LogLevel:       os.Getenv("LOG_LEVEL"),
		LogDevelopMode: os.Getenv("LOG_DEVELOP_MODE") == "true",
		AppName:        os.Getenv("APP_NAME"),
		AppVersion:     os.Getenv("APP_VERSION"),
		AllowedOrigins: strings.Split(os.Getenv("CORS_ORIGINS"), ","),
	}
	if db, err := LoadDbConfig(); err != nil {
		return &Config{}, err
	} else {
		cfg.DB = db
	}
	if server, err := LoadServerConfig(); err != nil {
		return &Config{}, err
	} else {
		cfg.Server = server
	}
	return &cfg, nil
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

func (db DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Pass, db.Name,
	)
}
