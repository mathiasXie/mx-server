package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/pkg/metrics"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime).Seconds()

		// 更新指标
		status := c.Writer.Status()
		metrics.HttpRequestsTotal.WithLabelValues(c.FullPath(), c.Request.Method, http.StatusText(status)).Inc()
		metrics.HttpRequestDuration.WithLabelValues(c.FullPath(), c.Request.Method).Observe(duration)
	}
}
