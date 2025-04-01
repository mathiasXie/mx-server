package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/gin-web/dto"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/factory"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/service"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

type SysHandler struct {
	Service *service.SysService
}

func (h *SysHandler) RouteArco(ctx *gin.Context) {
	val, exists := ctx.Get("userID")
	if !exists {
		dto.LoginStatusExceptionResp(ctx)
		return
	}
	userId, ok := val.(uint64)
	if !ok {
		dto.LoginStatusExceptionResp(ctx)
		return
	}
	sysUser, _, userRoutes, _ := h.Service.FindOneUser(ctx, int64(userId))
	var userRouteIds []int64
	if sysUser.IsSuper == 1 {
		routes, _ := h.Service.GetAllSysRoute(ctx)
		for _, v := range routes {
			userRouteIds = append(userRouteIds, v.Id)
		}
	} else {
		for _, v := range userRoutes {
			userRouteIds = append(userRouteIds, v.Id)
		}
	}
	routes := h.Routes(ctx, 0, userRouteIds)
	dto.Success(ctx, routes)
}

func (h *SysHandler) Routes(ctx *gin.Context, parentId int64, userRouteIds []int64) []dto.ArcoRouteResponse {

	result, _ := h.Service.BuildMenuRoute(ctx, userRouteIds, parentId)
	var routes []dto.ArcoRouteResponse
	for _, v := range result {
		route := dto.ArcoRouteResponse{
			Name:       v.Label,
			Key:        v.Path,
			Icon:       v.Icon,
			BreadCrumb: true,
			Children:   h.Routes(ctx, v.Id, userRouteIds),
		}
		routes = append(routes, route)
	}
	return routes
}

func (h *SysHandler) Login(ctx *gin.Context) {
	req := dto.LoginRequest{}
	factory.RequestCheck(ctx, &req)
	//if !utils.VerifyCaptcha(req.VerifyKey, req.VerifyCode, loader.GetRedis()) {
	//	dto.Fail(ctx, http.StatusBadRequest, "verify code error")
	//	return
	//}

	resp, err := h.Service.Login(ctx, &req)
	if err != nil {
		logger.CtxError(ctx, err.Error())
		dto.InternalErrorResp(ctx)
		return
	}
	if resp == nil {
		dto.Fail(ctx, http.StatusForbidden, "请检查用户名和密码")
	} else {
		dto.Success(ctx, resp)
	}
}

func (h *SysHandler) GetUserInfo(ctx *gin.Context) {
	val, exists := ctx.Get("userID")
	if !exists {
		dto.LoginStatusExceptionResp(ctx)
		return
	}
	userId, ok := val.(uint64)
	if !ok {
		dto.LoginStatusExceptionResp(ctx)
		return
	}
	sysUser, userRole, userRoutes, err := h.Service.FindOneUser(ctx, int64(userId))
	if err != nil {
		dto.LoginStatusExceptionResp(ctx)
		return
	}
	var userRoleRoutes []string
	if sysUser.IsSuper == 1 {
		userRoleRoutes = append(userRoleRoutes, "super")
	} else {
		for _, v := range userRoutes {
			userRoleRoutes = append(userRoleRoutes, v.Route.Name)
		}
	}
	dto.Success(ctx, dto.GetUserInfoResponse{
		UserId:   int64(userId),
		UserName: sysUser.Username,
		Role:     userRole.Name,
		UserRole: userRoleRoutes,
	})
}

func (h *SysHandler) RouteList(ctx *gin.Context) {

	dto.Success(ctx, h.sysRoutes(ctx, 0))
}

func (h *SysHandler) sysRoutes(ctx *gin.Context, parentId int64) []dto.RouteListResponse {
	result, _ := h.Service.GetRouteByParentID(ctx, parentId)
	var routes []dto.RouteListResponse
	for _, v := range result {
		route := dto.RouteListResponse{
			ID:        v.Id,
			Key:       v.Id,
			Name:      v.Name,
			Label:     v.Label,
			ParentID:  v.ParentId,
			Path:      v.Path,
			ApiPath:   v.ApiPath,
			Icon:      v.Icon,
			Sequence:  v.Sequence,
			Type:      int32(v.Type),
			Status:    int32(v.Status),
			Component: v.Component,
			Children:  h.sysRoutes(ctx, v.Id),
		}
		routes = append(routes, route)
	}
	return routes
}

func (h *SysHandler) UserList(ctx *gin.Context) {
	resp, err := h.Service.GetUser(ctx)
	if err != nil {
		logger.CtxError(ctx, "[SysHandler]UserList error:", err.Error())
		dto.InternalErrorResp(ctx)
		return
	}
	dto.Success(ctx, resp)
}
func (h *SysHandler) UserUpdate(ctx *gin.Context) {
	req := dto.UserUpdateRequest{}
	factory.RequestCheck(ctx, &req)
	err := h.Service.UpdateUser(ctx, &req)
	if err != nil {
		logger.CtxError(ctx, "[SysHandler]UserUpdate error:", err.Error())
		dto.InternalErrorResp(ctx)
		return
	}
	dto.Success(ctx, "ok")
}
func (h *SysHandler) UserDelete(ctx *gin.Context) {

	dto.Success(ctx, "ok")
}
func (h *SysHandler) UserResetPassword(ctx *gin.Context) {
	dto.Success(ctx, "ok")
}

func (h *SysHandler) RouteUpdate(ctx *gin.Context) {
	req := dto.RouteUpdateRequest{}
	reqCheck := factory.RequestCheck(ctx, &req)
	if reqCheck != "" {
		dto.Fail(ctx, http.StatusBadRequest, reqCheck)
		return
	}
	err := h.Service.RouteUpdate(ctx, req)
	if err != nil {
		dto.InternalErrorResp(ctx)
	}
	dto.Success(ctx, "ok")
}

func (h *SysHandler) RouteDelete(ctx *gin.Context) {
	req := dto.RouteDeleteRequest{}
	reqCheck := factory.RequestCheck(ctx, &req)
	if reqCheck != "" {
		dto.Fail(ctx, http.StatusBadRequest, reqCheck)
		return
	}
	err := h.Service.RouteDelete(ctx, req)
	if err != nil {
		dto.InternalErrorResp(ctx)
	}

}
