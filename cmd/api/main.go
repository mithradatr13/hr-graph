package main

import (
	"log/slog"

	"task-manager/internal/domain"
	"task-manager/internal/handler"
	"task-manager/internal/middleware"
	"task-manager/internal/repository"
	"task-manager/internal/service"
	"task-manager/pkg/cache"
	"task-manager/pkg/config"
	"task-manager/pkg/database"

	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task-manager/internal/router"
	"task-manager/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		baseLogger.Info("Starting server", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			baseLogger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":9090", mux)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	baseLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		baseLogger.Error("Server forced to shutdown", "error", err)
	}

	baseLogger.Info("Server exited properly")
}
