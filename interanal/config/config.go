package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct { //todo
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEV" env-default:"false"`
	Listen        struct {
		Type string `env:"LISTEN_TYPE" env-default:"port"`
		Ip   string `env:"IP" env-default:"0.0.0.0"`
		Port string `env:"PORT" env-default:"8000"`
	}
	AppConfig struct {
		LogLevel  string `env:"LOG_LEVEL" env-default:"trace"`
		AdminUser struct {
			Email    string `env:"ADMIN_EMAIL" env-default:"admin"`
			Password string `env:"ADMIN_PWD" env-default:"admin"`
		}
	}
	PostgreSQL struct {
		Username string `env:"DB_USER" env-required:"true"`
		Password string `env:"DB_PASSWORD" env-required:"true"`
		Host     string `env:"PSQL_HOST" env-required:"true"`
		Port     string `env:"PSQL_PORT" env-required:"true"`
		Database string `env:"DB_NAME" env-required:"true"`
	}
}

//const (
//	EnvConfigPathName  = "CONFIG-PATH"
//	FlagConfigPathName = "config"
//)

// var configPath string
var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		//flag.StringVar(&configPath, "config", ".env", "this is app config file")
		//flag.Parse()

		log.Print("config init")

		//if configPath == "" {
		//	configPath = os.Getenv(EnvConfigPathName)
		//}
		//
		//if configPath == "" {
		//	log.Fatal("config path is required")
		//}

		instance = &Config{}

		if err := cleanenv.ReadConfig(".env", instance); err != nil {
			text := "The Art of Development - Production Service"
			description, _ := cleanenv.GetDescription(instance, &text)
			log.Print(description)
			log.Fatal(err)
		}
	})

	return instance
}
