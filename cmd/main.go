package main

import (
	"github.com/sirupsen/logrus"
	"userSegments/interanal/app"
	"userSegments/interanal/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		logrus.WithError(err).Fatal("config parse error")
	}

	a, err := app.NewApp(cfg)
	if err != nil {
		logrus.WithError(err).Fatal("app not started")
	}

	a.Run()
}
