package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"online-judge/internal/config"
	"online-judge/internal/database"
	"online-judge/internal/judge"
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

	if err := os.MkdirAll(cfg.Judge.WorkDir, 0755); err != nil {
		logger.Fatal("create work dir", zap.Error(err))
	}

	problemRepo := repository.NewProblemRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)
	testCaseRepo := repository.NewTestCaseRepository(db)

	engine := judge.NewEngine(logger, cfg.Judge.WorkDir, cfg.Judge.CompileTimeout)
	worker := judge.NewWorker(engine, submissionRepo, problemRepo, testCaseRepo, logger)

	consumer, err := mq.NewConsumer(&cfg.RocketMQ, worker.Handle)
	if err != nil {
		logger.Fatal("create rocketmq consumer", zap.Error(err))
	}
	if err := consumer.Start(); err != nil {
		logger.Fatal("start rocketmq consumer", zap.Error(err))
	}
	defer consumer.Close()

	logger.Info("judge worker started",
		zap.String("topic", cfg.RocketMQ.Topic),
		zap.String("group", cfg.RocketMQ.ConsumerGroup),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("worker shutting down")
}
