package middleware

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/gin-web/dto"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/mathiasXie/gin-web/pkg/metrics"
)

func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				logger.CtxPanic(c, err, string(stack))
				dto.InternalErrorResp(c)
				metrics.ApplicationPanicTotal.WithLabelValues(string(stack)).Inc()
			}
		}()
		c.Next()
	}
}
