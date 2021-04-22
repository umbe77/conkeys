package config

import (
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type MongoDb struct {
	ConnectionUri string `envconfig:"mongo_connectionuri"`
}

type Config struct {
	Provider string `envconfig:"provider"`
	Mongo MongoDb
}

var cfg Config

var doOnce sync.Once

func GetConfig() *Config {
	doOnce.Do(func() {
		cfg = Config{
			Provider: "memory",
		}

		err := envconfig.Process("", &cfg)
		if err != nil {
			os.Exit(2)
		}
	})
	return &cfg
}
