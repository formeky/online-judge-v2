package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"online-judge/internal/api"
	"online-judge/internal/api/handler"
	"online-judge/internal/api/middleware"
	"online-judge/internal/config"
	"online-judge/internal/database"
	"online-judge/internal/mq"
	"online-judge/internal/repository"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := database.NewMySQL(&cfg.Database)
	if err != nil {
		logger.Fatal("connect mysql", zap.Error(err))
	}
	if err := database.AutoMigrate(db); err != nil {
		logger.Fatal("auto migrate", zap.Error(err))
	}

	producer, err := mq.NewProducer(&cfg.RocketMQ)
	if err != nil {
		logger.Fatal("create rocketmq producer", zap.Error(err))
	}
	defer producer.Close()

	problemRepo := repository.NewProblemRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)
	testCaseRepo := repository.NewTestCaseRepository(db)

	middleware.SetJWTSecret(cfg.JWT.Secret)
	gin.SetMode(cfg.Server.Mode)

	problemHandler := handler.NewProblemHandler(problemRepo, testCaseRepo)
	submissionHandler := handler.NewSubmissionHandler(submissionRepo, problemRepo, producer)

	router := api.SetupRouter(logger, problemHandler, submissionHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		logger.Info("server starting", zap.String("addr", addr))
		if err := router.Run(addr); err != nil {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("server shutting down")
}
