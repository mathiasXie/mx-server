package router

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mathiasXie/gin-web/applications/tts-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/handler"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/loader/resource"
	"github.com/mathiasXie/gin-web/consts"
	"google.golang.org/grpc/metadata"
)

func InitRouter(ctx context.Context, r *gin.Engine) *gin.Engine {

	r.POST("/tts", func(ctx *gin.Context) {
		text := ctx.PostForm("text")

		//简单进行一下测试
		client := resource.GetResource().TTSRpcClient.TTSServiceClient

		// 创建带有元数据的上下文
		trace_id, _ := ctx.Value(consts.LogID).(string)
		md := metadata.Pairs("trace_id", trace_id)
		rpcCtx := metadata.NewOutgoingContext(ctx, md)

		resp, err := client.TextToSpeech(rpcCtx, &proto.TextToSpeechRequest{
			Provider: proto.Provider_VOLCENGINE,
			VoiceId:  "zh_female_wanwanxiaohe_moon_bigtts",
			Language: "zh-CN",
			Text:     text,
		})
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": err.Error(),
			})
		} else {

			// 将响应写入文件
			filePath := "./tts_response.mp3"
			file, err := os.Create(filePath)
			if err != nil {
				ctx.JSON(500, gin.H{
					"message": "failed to create file",
				})
			}
			defer file.Close()
			_, err = file.Write(resp.AudioData)
			if err != nil {
				ctx.JSON(500, gin.H{
					"message": "failed to write file",
				})
			}
			ctx.JSON(200, map[string]interface{}{
				"message": "success",
				"data":    "ok",
			})
		}
	})

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
	r.GET("/xiaozhi/v1/", func(ctx *gin.Context) {
		// 创建带有元数据的上下文
		chatHandler := handler.ChatHandler{
			LLMClient: &llmClient,
			TTSClient: &ttsClient,
		}
		chatHandler.Chat(ctx, upgrader)
	})
	return r
}
