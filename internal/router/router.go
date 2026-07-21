package router

import (
	"net/http"
	"net/http/pprof"
	"task-manager/internal/handler"
	"task-manager/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(taskHandler *handler.TaskHandler) *gin.Engine {
	r := gin.Default()

	// اتصال میدلور مانیتورینگ پرومتئوس و ترسینگ لایو
	r.Use(middleware.MetricsMiddleware())

	// سرویس مستندات آنلاین Swagger UI به همراه فایل OpenAPI
	r.StaticFile("/docs/openapi.yaml", "./docs/openapi.yaml")
	r.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4/swagger-ui.css">
				<script src="https://unpkg.com/swagger-ui-dist@4/swagger-ui-bundle.js"></script>
			</head>
			<body>
				<div id="swagger-ui"></div>
				<script>
					SwaggerUIBundle({ url: '/docs/openapi.yaml', dom_id: '#swagger-ui' });
				</script>
			</body>
			</html>
		`))
	})

	// اندپوینت متریک‌های سیستم برای ابزار Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// ثبت سیستم Profiling (pprof) جهت تحلیل پرفورمنس و بنچمارک لود تست
	pprofGroup := r.Group("/debug/pprof")
	{
		pprofGroup.GET("/", gin.WrapH(http.HandlerFunc(pprof.Index)))
		pprofGroup.GET("/profile", gin.WrapH(http.HandlerFunc(pprof.Profile)))
		pprofGroup.GET("/cmdline", gin.WrapH(http.HandlerFunc(pprof.Cmdline)))
		pprofGroup.GET("/symbol", gin.WrapH(http.HandlerFunc(pprof.Symbol)))
		pprofGroup.GET("/trace", gin.WrapH(http.HandlerFunc(pprof.Trace)))
	}

	// v1 := r.Group("/api/v1/tasks")
	// {
	r.POST("/api/v1/tasks", taskHandler.Create)
	r.GET("/api/v1/tasks", taskHandler.List)
	r.GET("/api/v1/tasks/:id", taskHandler.GetByID)
	r.PUT("/api/v1/tasks/:id", taskHandler.Update)
	r.DELETE("/api/v1/tasks/:id", taskHandler.Delete)
	// }

	return r
}
