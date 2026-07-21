package main

import (
	"log"
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

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	log.Println("running Task Manager Service...")

	cfg := config.Load()

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("❌ Connection to PostgreSQL failed: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("❌ Failed to create database tables: %v", err)
	}

	rdb, err := cache.NewRedisClient(cfg)
	if err != nil {
		log.Printf("⚠️ Warning: Failed to connect to Redis. System will operate without caching: %v", err)
	}

	repo := repository.NewPostgresTaskRepository(db)

	var taskCache domain.TaskCache = nil
	if rdb != nil {
		taskCache = repository.NewRedisTaskCache(rdb, 10*time.Minute)
	}

	svc := service.NewTaskService(repo, taskCache)
	taskHandler := handler.NewTaskHandler(svc)

	prometheus.MustRegister(service.TasksCountGauge)
	middleware.RegisterMetrics()

	r := router.SetupRouter(taskHandler)

	log.Printf("🛰️ Service is running on port :%s.", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("❌ Failed to run server: %v", err)
	}
}
