package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var LogLevels = map[string]int{
	"DEBUG": -4,
	"INFO":  0,
	"WARN":  1,
	"ERROR": 8,
}

type Config struct {
	ServiceCnt int `env:"SERVICE_CNT"`
	Port       int

	Logger
}

type Logger struct {
	Level string `env:"LOG_LEVEL"`
}

func GetConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		fmt.Printf("environment is not OK: %s\n", err)
		os.Exit(1)
	}

	return &cfg
}

func GetConfigRun() *Config {
	return &Config{
		ServiceCnt: 5,
		Port:       8090,
		Logger: Logger{
			Level: "DEBUG",
		},
	}
}
