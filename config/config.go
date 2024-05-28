package config

import (
	"os"
	"time"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string 			`yaml:"env" env-default:"local"`
	StorageUrl string 	`yaml:"storage_url" env-required:"true"`
	HTTPServer			`yaml:"http_server"`
}

type HTTPServer struct {
	Address 	string 			`yaml:"address" env-default:"localhost:8089"`
	Timeout 	time.Duration 	`yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration 	`yaml:"idle-timeout" env-default:"60s"`
}

func MustLoad(cfgPath string) *Config {
	if _, err := os.Stat(cfgPath); err != nil {
		log.Fatalf("config file does not exist: %s", cfgPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		log.Fatalf("canno`t read config file: %s", cfgPath)
	}

	return &cfg
}