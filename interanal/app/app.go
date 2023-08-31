package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
	"userSegments/interanal/config"
	"userSegments/interanal/controller"
	"userSegments/interanal/service"
	"userSegments/interanal/storage"
	psql "userSegments/pkg/postgresql"
)

type App struct {
	cfg                 *config.Config
	router              *httprouter.Router
	httpServer          *http.Server
	logger              *logrus.Entry
	userService         service.User
	segmentService      service.Segment
	userSegmentsService service.UserSegments
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

	//metricHandler := metric.Handler{}
	//metricHandler.Register(router)
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
		return App{}, errors.New("psql.NewClient")
	}

	// TODO вынести роуты
	userStorage := storage.NewUserStorage(pgClient)
	segmentStorage := storage.NewSegmentStorage(pgClient)
	userSegmentsStorage := storage.NewUserSegmentsStorage(pgClient)
	userService := service.NewUserService(&userStorage)
	segmentService := service.NewSegmentService(&segmentStorage)
	userSegmentsService := service.NewUserSegmentsService(&userSegmentsStorage, segmentStorage)
	userController := controller.NewUserController(userService)
	segmentController := controller.NewSegmentController(segmentService)
	usersSegmentController := controller.NewUserSegmentsController(userSegmentsService)

	router.GET("/api/users/:id", userController.GetUser)
	router.POST("/api/users", userController.CreateUser)

	router.GET("/api/segments/:id", segmentController.GetSegment) //todo :slug ????
	router.POST("/api/segments", segmentController.CreateSegment)
	router.DELETE("/api/segments/:slug", segmentController.DeleteSegment)

	router.GET("/api/users/:id/segments", usersSegmentController.GetUserSegments)
	router.POST("/api/users/:id/segments", usersSegmentController.AddUserToSegment)

	return App{
		cfg:                 cfg,
		router:              router,
		logger:              logger,
		userService:         userService,
		segmentService:      segmentService,
		userSegmentsService: userSegmentsService,
	}, nil
}

func (a *App) Run() {
	a.startHTTP()
}

func (a *App) startHTTP() {
	logrus.Info("HTTP Server initializing")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.Ip, a.cfg.Listen.Port))
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
}
