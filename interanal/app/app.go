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

func NewApp(cfg *config.Config) (App, error) {
	router := httprouter.New()
	logger := e

	pgConfig := psql.NewPgConfig(
		cfg.PostgreSQL.Username,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Database,
	)

	pgClient, err := psql.NewClient(context.Background(), 3, 3*time.Second, pgConfig)
	if err != nil {
		return App{}, errors.New("psql client error")
	}

	userStorage := storage.NewUserStorage(pgClient)
	segmentStorage := storage.NewSegmentStorage(pgClient)
	userSegmentsStorage := storage.NewUserSegmentsStorage(pgClient)
	segmentHistoryStorage := storage.NewSegmentHistoryStorage(pgClient)

	userService := service.NewUserService(&userStorage)
	segmentService := service.NewSegmentService(&segmentStorage, &userSegmentsStorage)
	userSegmentsService := service.NewUserSegmentsService(&userSegmentsStorage, segmentStorage)
	segmentHistoryService := service.NewSegmentsHistoryService(&segmentHistoryStorage)

	userController := controller.NewUserController(userService)
	segmentController := controller.NewSegmentController(segmentService)
	usersSegmentController := controller.NewUserSegmentsController(userSegmentsService)
	segmentHistoryController := controller.NewSegmentHistoryControllerController(segmentHistoryService)

	router.GET("/api/users/:id", userController.GetUser)
	router.POST("/api/users", userController.CreateUser)

	router.GET("/api/segments/:id", segmentController.GetSegment)
	router.POST("/api/segments", segmentController.CreateSegment)
	router.DELETE("/api/segments/:slug", segmentController.DeleteSegment)

	router.GET("/api/users/:id/segments", usersSegmentController.GetUserSegments)
	router.POST("/api/users/:id/segments", usersSegmentController.ChangeUserSegments)

	router.POST("/api/segment-history/generate", segmentHistoryController.GenerateHistoryReport)
	router.GET("/reports/:filename", segmentHistoryController.DownloadReport)

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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", a.cfg.Listen.Port))
	if err != nil {
		a.logger.Fatal("failed to create listener")
	}

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
}
