package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	"github.com/mathiasXie/gin-web/utils"
)

// AuthRoleMiddleware 权限验证中间件
func AuthRoleMiddleware(ctx *gin.Context) {

	token := ctx.Request.Header.Get("Access-Token")
	if token == "" {
		// 处理 没有token的时候
		dto.Fail(ctx, http.StatusForbidden, "登录失效，请重新登录")
		ctx.Abort() // 不会继续停止
		return
	}
	// 解析
	mc, err := utils.ParseToken(token)
	if err != nil {
		// 处理 解析失败
		dto.Fail(ctx, http.StatusForbidden, "登录信息解析失败，请重新登录")
		ctx.Abort()
		return
	}

	// 将当前请求的userID信息保存到请求的上下文c上
	ctx.Set("userID", mc.UserID)
	ctx.Next()
}
