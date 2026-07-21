package main

import (
	"log/slog"
	"time"

	"task-manager/internal/domain"
	"task-manager/internal/handler"
	"task-manager/internal/middleware"
	"task-manager/internal/repository"
	"task-manager/internal/router"
	"task-manager/internal/service"
	"task-manager/pkg/cache"
	"task-manager/pkg/config"
	"task-manager/pkg/database"
	"task-manager/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	slog.Info("🚀 Starting Task Manager Service...")

	cfg := config.Load()

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		slog.Error("❌ Failed to connect to PostgreSQL", "error", err)
		return
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		slog.Error("❌ Failed to create database tables: %v", err)
		return
	}

	rdb, err := cache.NewRedisClient(cfg)
	if err != nil {
		slog.Warn("⚠️ Warning: Failed to connect to Redis. System will operate without caching: %v", err)
	}

	repo := repository.NewPostgresTaskRepository(db)

	var taskCache domain.TaskCache = nil
	if rdb != nil {
		taskCache = repository.NewRedisTaskCache(rdb, 10*time.Minute)
	}

	baseLogger := logger.InitLogger()
	svc := service.NewTaskService(repo, taskCache, baseLogger)
	taskHandler := handler.NewTaskHandler(svc)

	prometheus.MustRegister(service.TasksCountGauge)
	middleware.RegisterMetrics()

	r := router.SetupRouter(taskHandler)

	slog.Info("🛰️ Service is running on port :%s.", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		slog.Error("❌ Failed to run server: %v", err)
	}
}
