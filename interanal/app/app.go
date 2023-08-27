package app

import (
	"context"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
	"userSegments/interanal/config"
	"userSegments/interanal/controller"
	"userSegments/interanal/domain/user/service"
	"userSegments/interanal/domain/user/storage"
	"userSegments/pkg/metric"
	psql "userSegments/pkg/postgresql"
)

type App struct {
	cfg        *config.Config
	router     *httprouter.Router
	httpServer *http.Server
	//logger     *Logger
	logger *logrus.Entry
}

var e *logrus.Entry

//type Logger struct {
//	*logrus.Entry
//}

//func GetLogger() *Logger {
//	return &Logger{e}
//}

func NewApp(cfg *config.Config) (App, error) {
	router := httprouter.New()
	logger := e

	//router.GET("/", Index)
	router.GET("/test", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		logrus.Print("test")
		w.WriteHeader(204)
		w.Write([]byte("test"))
	})

	//router.HandlerFunc(http.MethodGet, "/test", func(w http.ResponseWriter, r *http.Request) {
	//	logrus.Print("test")
	//	w.WriteHeader(204)
	//	w.Write([]byte("test"))
	//})
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	//log.Fatal(http.ListenAndServe(":3333", router))

	//pgDsn := fmt.Sprintf(
	//	"postgres://%s:%s@%s:%s/%s",
	//	cfg.PostgreSQL.Username,
	//	cfg.PostgreSQL.Password,
	//	cfg.PostgreSQL.Host,
	//	cfg.PostgreSQL.Port,
	//	cfg.PostgreSQL.Database,
	//)
	pgConfig := psql.NewPgConfig(
		cfg.PostgreSQL.Username,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Database,
	)
	//pgConfig.ConnStringFromCfg()

	// TODO to config
	pgClient, err := psql.NewClient(context.Background(), 5, 3*time.Second, pgConfig)
	if err != nil {
		//return App{}, errors.Wrap(err, "psql.NewClient")
	}

	//closer.AddN(pgClient)

	userStorage := storage.NewUserStorage(pgClient)
	//user, err := userStorage.GetById(context.Background(), 1) // todo 1
	//if err != nil {
	//	logrus.Fatal(err)
	//}
	//logrus.Fatal(user)
	//productStorage := storage.NewUserStorage(pgClient)
	userService := service.NewUserService(&userStorage)
	//user, _ := userService.GetUserById(context.Background(), 1)
	//if err != nil {
	//	logrus.Fatal(err)
	//}
	//logrus.Fatal(user)
	userController := controller.NewUserController(userService)

	router.HandlerFunc(http.MethodGet, "/api/user", userController.GetUser)

	return App{
		cfg:    cfg,
		router: router,
		logger: logger,
		//productServiceServer: userStorage,
	}, nil
}

func (a *App) Run() {
	a.startHTTP()
}

func (a *App) startHTTP() {
	//logger := logging.WithFields(ctx,
	//	logging.StringField("IP", a.cfg.HTTP.IP),
	//	logging.IntField("Port", a.cfg.HTTP.Port),
	//)
	logrus.Info("HTTP Server initializing")

	//listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "", 3333)) // todo
	listener, err := net.Listen("tcp", ":8000") // todo cfg.PostgreSQL.Port,
	//listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.Listen.Ip, a.cfg.Listen.Port))
	if err != nil {
		a.logger.Fatal("failed to create listener")
	}

	//handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      a.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logrus.Warn("server shutdown")
		default:
			logrus.Fatal("failed to start server")
		}
	}

	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		a.logger.Fatal("failed to shutdown server")
	}

	//return err
}
