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
	log.Println("🚀 در حال راه‌اندازی میکروسرویس مدیریت تسک...")

	cfg := config.Load()

	// ۲. اتصال لایه دیتابیس پایدار Postgres
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("❌ اتصال به PostgreSQL با خطا مواجه شد: %v", err)
	}
	defer db.Close()

	// اجرای خودکار مهاجرت دیتابیس برای ساخت تیبل‌ها
	if err := database.Migrate(db); err != nil {
		log.Fatalf("❌ ایجاد جدول دیتابیس ناموفق بود: %v", err)
	}

	// ۳. اتصال کلاینت حافظه موقت Redis
	rdb, err := cache.NewRedisClient(cfg)
	if err != nil {
		log.Printf("⚠️ هشدار: اتصال به Redis برقرار نشد. سیستم بدون کش ادامه می‌دهد: %v", err)
	}

	// ۴. تزریق وابستگی‌ها (Dependency Injection)
	repo := repository.NewPostgresTaskRepository(db)

	var taskCache domain.TaskCache = nil
	if rdb != nil {
		taskCache = repository.NewRedisTaskCache(rdb, 10*time.Minute)
	}

	svc := service.NewTaskService(repo, taskCache)
	taskHandler := handler.NewTaskHandler(svc)

	// ثبت متریک‌های سفارشی در سیستم Prometheus
	prometheus.MustRegister(service.TasksCountGauge)
	middleware.RegisterMetrics()

	// ۵. مقداردهی روتر و راه‌اندازی پورت سرور
	r := router.SetupRouter(taskHandler)

	log.Printf("🛰️ سرویس با موفقیت روی پورت :%s در دسترس قرار گرفت.", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("❌ اجرای سرور با شکست مواجه شد: %v", err)
	}
}
