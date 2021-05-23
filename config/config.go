package config

import (
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type MongoDb struct {
	ConnectionUri string `envconfig:"mongo_connectionuri"`
}

type PostegresSQL struct {
	ConnectionUri string `envconfig:"postgres_connectionuri"`
}

type Admin struct {
	Password string `envconfig:"user_password"`
}

type Config struct {
	Provider string `envconfig:"provider"`
	Admin    Admin
	Mongo    MongoDb
	Postgres PostegresSQL
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
