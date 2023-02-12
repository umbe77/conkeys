package config

import (
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type PostegresSQL struct {
	ConnectionUri string `envconfig:"postgres_connectionuri"`
}

type Admin struct {
	Password string `envconfig:"admin_password"`
}

type Config struct {
	Admin    Admin
	Postgres PostegresSQL
}

var cfg Config

var doOnce sync.Once

func GetConfig() *Config {
	doOnce.Do(func() {
		cfg = Config{}

		err := envconfig.Process("", &cfg)
		if err != nil {
			os.Exit(2)
		}
	})
	return &cfg
}
