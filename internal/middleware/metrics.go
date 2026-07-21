package middleware

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of HTTP requests processed by the service",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_latency_histogram",
			Help:    "Latency of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpLatencyHistogram)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = strconv.FormatInt(time.Now().UnixNano(), 16)
		}
		c.Header("X-Trace-ID", traceID)

		c.Next()

		latency := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		HttpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		HttpLatencyHistogram.WithLabelValues(c.Request.Method, path).Observe(latency)

		log.Printf("[TraceID: %s] %s %s | Status: %s | Latency: %v seconds", traceID, c.Request.Method, path, status, latency)
	}
}
