package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "X-Request-ID"

func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		reqID := c.GetHeader(RequestIDKey)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set(RequestIDKey, reqID)
		c.Header(RequestIDKey, reqID)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		slog.Info("HTTP Request",
			slog.String("request_id", reqID),
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", rawQuery),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
			slog.String("user_agent", c.Request.UserAgent()),
		)
	}
}
