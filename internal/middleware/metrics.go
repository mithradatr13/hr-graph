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
			Help: "تعداد کل درخواست‌های پردازش شده توسط سرویس HTTP",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_latency_histogram",
			Help:    "میزان تاخیر درخواست‌ها بر حسب ثانیه",
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

		// مکانیزم Tracing ساده: تزریق شناسه یکتا به هدر پاسخ برای رهگیری عملیات کلاینت
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
		
		log.Printf("[TraceID: %s] %s %s | وضعیت: %s | تاخیر: %v ثانیه", traceID, c.Request.Method, path, status, latency)
	}
}
