package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/gin-web/dto"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/factory"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/service"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

type UserHandler struct {
	Service *service.UserService
}

func (h *UserHandler) AddUser(ctx *gin.Context) {
	req := &dto.AddUserRequest{}
	factory.RequestCheck(ctx, &req)

	err := h.Service.AddUser(ctx, req)
	if err != nil {
		logger.CtxError(ctx, "[UserHandler] AddUser Error", err.Error())
		dto.InternalErrorResp(ctx)
		return
	}
	dto.Success(ctx, "ok")
}

func (h *UserHandler) GetUsers(ctx *gin.Context) {
	users, total, err := h.Service.GetUsers(ctx)
	if err != nil {
		logger.CtxError(ctx, "[UserHandler]GetUsers Error:", err.Error())
		dto.InternalErrorResp(ctx)
		return
	}
	time.Sleep(time.Second * 3)
	dto.Success(ctx, &dto.GetAllUserResponse{
		UserList: users,
		Total:    total,
	})
}
