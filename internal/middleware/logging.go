package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nevzattalhaozcan/forgotten/pkg/logger"
	"go.uber.org/zap"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		method := c.Request.Method

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		}

		if raw != "" {
			fields = append(fields, zap.String("query", raw))
		}

		if statusCode >= 500 {
			logger.Error("Server error", fields...)
		} else if statusCode >= 400 {
			logger.Warn("Client error", fields...)
		} else {
			logger.Info("Request processed", fields...)
		}
	}
}