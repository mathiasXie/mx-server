package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/gin-web/dto"
	"golang.org/x/time/rate"
)

// RateLimit 限流中间件
func RateLimit(limit, burst int) gin.HandlerFunc {

	limiter := rate.NewLimiter(rate.Limit(limit), burst)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			dto.Fail(c, http.StatusForbidden, "服务器繁忙，请稍候再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
