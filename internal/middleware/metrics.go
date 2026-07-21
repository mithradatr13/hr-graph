package middleware

import (
	"log/slog"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed by the service",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Latency of HTTP requests in seconds",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
		},
		[]string{"method", "endpoint"},
	)

	HttpActiveRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Current number of active HTTP requests being processed",
		},
	)

	GoroutinesCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_routines_current",
			Help: "Current number of goroutines running in the runtime",
		},
	)

	MemoryAllocBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_memory_alloc_bytes",
			Help: "Current bytes of allocated heap objects",
		},
	)

	MemorySysBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_memory_sys_bytes",
			Help: "Total bytes of memory obtained from the OS",
		},
	)

	GCPauseTotalNs = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "go_gc_pause_total_ns",
			Help: "Total GC pause time in nanoseconds",
		},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpLatencyHistogram)
	prometheus.MustRegister(HttpActiveRequests)
	prometheus.MustRegister(GoroutinesCount)
	prometheus.MustRegister(MemoryAllocBytes)
	prometheus.MustRegister(MemorySysBytes)
	prometheus.MustRegister(GCPauseTotalNs)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		HttpActiveRequests.Inc()
		defer HttpActiveRequests.Dec()

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

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		GoroutinesCount.Set(float64(runtime.NumGoroutine()))
		MemoryAllocBytes.Set(float64(m.Alloc))
		MemorySysBytes.Set(float64(m.Sys))
		GCPauseTotalNs.Add(float64(m.PauseTotalNs))

		slog.Info("HTTP request processed",
			"trace_id", traceID,
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"latency", latency,
			"goroutines", runtime.NumGoroutine(),
			"mem_alloc_bytes", m.Alloc,
		)
	}
}
