package logger

import (
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GinMiddleware returns a Gin middleware that injects a request-scoped logger into the request context
// and logs basic request information.
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// request id
		reqID := c.GetHeader(requestIDHeader)
		if reqID == "" {
			reqID = uuid.NewString()
			// also set header so downstream services/clients can see it
			c.Request.Header.Set(requestIDHeader, reqID)
		}

		// create request scoped logger
		l := slog.Default().With(
			slog.Group("client",
				slog.String("req_id", reqID),
				slog.String("method", c.Request.Method),
				slog.String("path", c.FullPath()),
				slog.String("remote", c.ClientIP()),
			),
		)

		// attach to context
		c.Request = c.Request.WithContext(NewContext(c.Request.Context(), l))

		// proceed
		c.Next()

		// after request
		duration := time.Since(start)
		status := c.Writer.Status()
		l.Info("request completed",
			slog.Int("status", status),
			slog.Duration("duration", duration),
			slog.String("handler", c.HandlerName()),
		)
	}
}
