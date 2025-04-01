package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/mathiasXie/gin-web/applications/gin-web/internal/model"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/service"
	"github.com/mathiasXie/gin-web/applications/gin-web/loader"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/gin-web/dto"
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

	// 进行权限验证
	if !checkPermission(ctx, int64(mc.UserID)) {
		dto.Fail(ctx, http.StatusForbidden, "登录信息解析失败，请重新登录")
		//http.Error(w, strconv.FormatInt(userId, 10), http.StatusForbidden) //无权限 403
		return
	}
	// 将当前请求的userID信息保存到请求的上下文c上
	ctx.Set("userID", mc.UserID)
	ctx.Next()
}
func checkPermission(ctx *gin.Context, userID int64) bool {

	cacheKey := "biz#sys#user_role_route"
	routeCache := loader.GetRedis().HGet(ctx, cacheKey, strconv.FormatInt(userID, 10))
	err := routeCache.Err()
	routeCacheContent := routeCache.Val()
	var userRoutes []model.SysRoute
	if err == nil && routeCacheContent != "" { //读到缓存，从缓存里取出信息
		if strings.TrimSpace(routeCacheContent) == "*" { //星号是超级管理员，直接通过
			return true
		}
		err := json.Unmarshal([]byte(routeCacheContent), &userRoutes)
		if err != nil {
			return false
		}
	} else { //未读到缓存，从数据库里取出信息并存入缓存
		sysService := service.NewSysService(ctx)
		sysUser, _, userRoutes, err := sysService.FindOneUser(ctx, userID)
		if err != nil {
			return false
		}

		var routeJson []byte
		if sysUser.IsSuper == 1 {
			routeJson = []byte("*")
		} else {
			routeJson, _ = json.Marshal(userRoutes)
		}
		err2 := loader.GetRedis().HSet(ctx, cacheKey, userID, string(routeJson))
		_ = loader.GetRedis().Expire(ctx, cacheKey, 3600*6)
		if err2 != nil {
			return false
		}
		if sysUser.IsSuper == 1 {
			return true
		}
	}

	parts := strings.Split(ctx.Request.URL.Path, "/")
	requestURI := "/" + strings.Join(parts[3:], "/")
	b := false
	for _, route := range userRoutes {
		if strings.Contains(route.ApiPath, requestURI) {
			b = true
			break
		}
	}
	return b
}
