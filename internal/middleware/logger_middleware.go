package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var Logger *slog.Logger

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		reqID, _ := c.Get(RequestIDKey)

		if raw != "" {
			path = path + "?" + raw
		}

		reqIDStr, _ := reqID.(string)

		attrs := []slog.Attr{
			slog.Int("status", statusCode),
			slog.Duration("latency", latency),
			slog.String("client_ip", clientIP),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("request_id", reqIDStr),
		}

		if statusCode >= 500 {
			Logger.LogAttrs(c.Request.Context(), slog.LevelError, "HTTP request failed", attrs...)
		} else if statusCode >= 400 {
			Logger.LogAttrs(c.Request.Context(), slog.LevelWarn, "HTTP request warning", attrs...)
		} else {
			Logger.LogAttrs(c.Request.Context(), slog.LevelInfo, "HTTP request success", attrs...)
		}
	}
}