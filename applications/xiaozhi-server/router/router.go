package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/handler"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/loader/resource"
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

	llmClient := resource.GetResource().LLMRpcClient.LLMServiceClient
	ttsClient := resource.GetResource().TTSRpcClient.TTSServiceClient
	asrClient := resource.GetResource().ASRRpcClient.ASRServiceClient
	r.GET("/xiaozhi/v1/", func(ctx *gin.Context) {
		// 创建带有元数据的上下文
		chatHandler := handler.ChatHandler{
			LLMClient: &llmClient,
			TTSClient: &ttsClient,
			ASRClient: &asrClient,
		}
		chatHandler.Chat(ctx, upgrader)
	})
	return r
}
