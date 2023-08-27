package main

import (
	"github.com/sirupsen/logrus"
	"userSegments/interanal/app"
	"userSegments/interanal/config"
)

func main() {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//logging.L(ctx).Info("config initializing")
	cfg := config.GetConfig()

	//ctx = logging.ContextWithLogger(ctx, logging.NewLogger())

	a, err := app.NewApp(cfg)
	if err != nil {
		logrus.WithError(err).Fatal("app.NewApp")
	}

	//logging.L(ctx).Info("Running Application")
	a.Run()
	//if err != nil {
	//	logging.WithError(ctx, err).Fatal("app.Run")
	//	return
	//}
}
