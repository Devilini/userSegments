package main

import (
	"github.com/sirupsen/logrus"
	"userSegments/interanal/app"
	"userSegments/interanal/config"
)

func main() {
	cfg := config.GetConfig()

	//ctx = logging.ContextWithLogger(ctx, logging.NewLogger())
	a, err := app.NewApp(cfg)
	if err != nil {
		logrus.WithError(err).Fatal("app.NewApp")
	}

	a.Run()
}
