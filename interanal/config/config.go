package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"sync"
)

type Config struct {
	IsDebug bool `env:"IS_DEBUG" env-default:"false"`
	Listen  struct {
		Ip   string `env:"IP" env-default:"0.0.0.0"`
		Port string `env:"PORT" env-default:"8000"`
	}
	AppConfig struct {
		LogLevel string `env:"LOG_LEVEL" env-default:"trace"`
	}
	PostgreSQL struct {
		Username string `env:"DB_USER" env-required:"true"`
		Password string `env:"DB_PASSWORD" env-required:"true"`
		Host     string `env:"PSQL_HOST" env-required:"true"`
		Port     string `env:"PSQL_PORT" env-required:"true"`
		Database string `env:"DB_NAME" env-required:"true"`
	}
}

var err error
var instance *Config
var once sync.Once

func GetConfig() (*Config, error) {
	once.Do(func() {
		logrus.Info("config init")

		instance = &Config{}
		err = cleanenv.ReadConfig(".env", instance)
		if err != nil {
			text := "Segment Service"
			description, _ := cleanenv.GetDescription(instance, &text)
			logrus.Info(description)
		}
	})

	return instance, err
}
