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

	// Attach Prometheus monitoring and live tracing middleware
	r.Use(middleware.MetricsMiddleware())

	// Serve online Swagger UI documentation alongside the OpenAPI spec file
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

	// System metrics endpoint for Prometheus scraping
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Register Profiling system (pprof) for performance analysis and load testing benchmarks
	pprofGroup := r.Group("/debug/pprof")
	{
		pprofGroup.GET("/", gin.WrapH(http.HandlerFunc(pprof.Index)))
		pprofGroup.GET("/profile", gin.WrapH(http.HandlerFunc(pprof.Profile)))
		pprofGroup.GET("/cmdline", gin.WrapH(http.HandlerFunc(pprof.Cmdline)))
		pprofGroup.GET("/symbol", gin.WrapH(http.HandlerFunc(pprof.Symbol)))
		pprofGroup.GET("/trace", gin.WrapH(http.HandlerFunc(pprof.Trace)))
	}

	// API v1 group and task resource routes
	v1 := r.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", taskHandler.Create)
			tasks.GET("", taskHandler.List)
			tasks.GET("/:id", taskHandler.GetByID)
			tasks.PUT("/:id", taskHandler.Update)
			tasks.DELETE("/:id", taskHandler.Delete)
		}
	}

	return r
}
