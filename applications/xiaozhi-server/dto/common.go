package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/consts"
)

var SuccessResponse = CommonResponse{
	Code: 200,
	Msg:  consts.SuccessMsg,
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &CommonResponse{
		Code: http.StatusOK,
		Msg:  consts.SuccessMsg,
		Data: data,
	})
}

func InternalErrorResp(ctx *gin.Context) {
	Fail(ctx, http.StatusInternalServerError, "服务内部错误")
}

func LoginStatusExceptionResp(ctx *gin.Context) {
	Fail(ctx, http.StatusInternalServerError, "登录状态问题")
}

func Fail(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, &CommonResponse{
		Code: code,
		Msg:  message,
	})
}

type CommonResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
