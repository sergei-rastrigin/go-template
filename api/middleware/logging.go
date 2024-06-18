package middleware

import (
	"context"
	"net/http"
	"time"

	appLogger "go-template/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func Logging(logger zerolog.Logger, skipPaths []string) gin.HandlerFunc {
	skipPathsMap := make(map[string]bool)
	for _, path := range skipPaths {
		skipPathsMap[path] = true
	}
	
	return func(c *gin.Context) {
		if _, ok := skipPathsMap[c.Request.URL.Path]; ok {
			c.Next()
			return
		}


		start := time.Now()

		traceID := c.Request.Header.Get("x-trace-id") // Case insensitive header
		if traceID == "" {
			traceID = uuid.NewString()
			logger.Debug().Str("trace_id", traceID).Msg("request received without \"x-trace-id\" - new trace id generated")
		}

		logger = logger.With().
			Str("trace_id", traceID).
			Logger()

		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		logger.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("Gin - Request received")

		// Add trace id and logger to context
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, appLogger.KeyTrace, traceID)
		ctx = context.WithValue(ctx, appLogger.KeyLog, logger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		responseLog := logger.With().
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("ip", c.ClientIP()).
			Str("latency", time.Since(start).String()).
			Str("user-agent", c.Request.UserAgent()).
			Logger()

		msg := "Gin - Request Ended"
		if len(c.Errors) > 0 {
			msg += ": " + c.Errors.String()
		}

		switch {
		case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
			responseLog.Warn().Msg(msg)
		case c.Writer.Status() >= http.StatusInternalServerError:
			responseLog.Error().Msg(msg)
		default:
			responseLog.Trace().Msg(msg)
		}
	}
}
