package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App
		HTTP
		Log
		PG
	}

	App struct {
		Name    string `env-required:"true" env:"APP_NAME"`
		Version string `env-required:"true" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT" env-default:"8080"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"info"`
	}

	PG struct {
		Host         string        `env-required:"true" env:"POSTGRES_HOST"`
		Port         string        `env-required:"true" env:"POSTGRES_PORT"`
		Username     string        `env-required:"true" env:"POSTGRES_USER"`
		Password     string        `env-required:"true" env:"POSTGRES_PASSWORD"`
		DBName       string        `env-required:"true" env:"POSTGRES_DB"`
		SSLMode      string        `env:"POSTGRES_SSL" env-default:"disable"`
		ConnAttempts int           `env:"CONN_ATTEMPTS" default:"10"`
		ConnTimeout  time.Duration `env:"CONN_TIMEOUT" default:"1s"`
	}
)

func NewConfig() *Config {

	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "path_to_config", ".env", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
