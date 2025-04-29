package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/handler"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/service"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/loader"
)

func InitRouter(ctx context.Context, r *gin.Engine) *gin.Engine {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许跨域请求，在生产环境中需要根据实际情况进行调整
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	r.GET("/xiaozhi/v1/", func(ctx *gin.Context) {
		chatHandler := handler.ChatHandler{
			LLMClient: loader.GetLLMRpc(),
			TTSClient: loader.GetTTSRpc(),
			ASRClient: loader.GetASRRpc(),
		}
		chatHandler.Chat(ctx, upgrader)
	})

	userGroup := r.Group("/user-center/v1")
	{
		userCenterHandler := handler.UserCenterHandler{
			UserService: service.NewUserService(ctx, loader.GetDB(ctx, true)),
		}
		userGroup.POST("/sign-in", func(ctx *gin.Context) {
			userCenterHandler.SignIn(ctx)
		})
		userGroup.POST("/sign-up", func(ctx *gin.Context) {
			userCenterHandler.SignUp(ctx)
		})
	}

	r.POST("/xiaozhi/ota/", func(ctx *gin.Context) {
		handler.OtaHandler(ctx)
	})
	// 激活状态检查
	r.POST("/xiaozhi/ota/activate", func(ctx *gin.Context) {
		handler.ActivateHandler(ctx)
	})
	return r
}
