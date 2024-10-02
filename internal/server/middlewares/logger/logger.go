package middleware_logger

import (
	"test-case/internal/utils/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		logger.Logger.Info().
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.String()).
			Str("client_ip", c.ClientIP()).
			Msg("Incoming request")

		c.Next()

		logger.Logger.Info().
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.String()).
			Str("client_ip", c.ClientIP()).
			Int("status_code", c.Writer.Status()).
			Dur("duration", time.Since(start)).
			Msg("Response")
	}
}
