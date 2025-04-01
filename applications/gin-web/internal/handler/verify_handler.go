package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/gin-web/dto"
	"github.com/mathiasXie/gin-web/applications/gin-web/loader"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/mathiasXie/gin-web/utils"
)

type VerifyHandler struct {
}

func (h *VerifyHandler) MakeVerify(ctx *gin.Context) {

	id, b64s, err := utils.GenerateCaptcha(ctx, loader.GetRedis())
	if err != nil {
		logger.CtxError(ctx, "[VerifyHandler] MakeVerify Error", err.Error())
		dto.InternalErrorResp(ctx)
		return
	}

	dto.Success(ctx, dto.VerifyData{
		VerifyKey:   id,
		VerifyImage: b64s,
	})
}
