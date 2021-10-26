package config

import (
	"os"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server *ServerConfig
	DB     *DBConfig
}

type ServerConfig struct {
	Port            uint16        `default:"8000"`
	IdleTimeout     time.Duration `default:"2m"`
	ReadTimeout     time.Duration `default:"5s"`
	WriteTimeout    time.Duration `default:"5s"`
	ShutdownTimeout time.Duration `default:"15s"`
}

type DBConfig struct {
	Host         string        `default:"localhost"`
	Port         uint16        `default:"5432"`
	Name         string        `default:"shortener"`
	User         string        `default:"shortener"`
	Password     string        `default:"123456789"`
	SSLMode      string        `envconfig:"ssl_mode" default:"verify-full"`
	QueryTimeout time.Duration `default:"500ms"`
}

func GetConfig() (*Config, error) {
	serverConf, err := parseServerConfig()
	if err != nil {
		return nil, err
	}

	dbConf, err := parseDBConfig()
	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port != "" {
		portUint, err := strconv.ParseUint(port, 10, 64)
		if err != nil {
			return nil, err
		}

		serverConf.Port = uint16(portUint)
	}

	return &Config{Server: serverConf, DB: dbConf}, nil
}

func GetDBConfig() (*DBConfig, error) {
	return parseDBConfig()
}

func parseServerConfig() (*ServerConfig, error) {
	var conf ServerConfig

	err := envconfig.Process("SERVER", &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func parseDBConfig() (*DBConfig, error) {
	var conf DBConfig

	err := envconfig.Process("DB", &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
