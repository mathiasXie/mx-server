package router

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/handler"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/service"
	"github.com/mathiasXie/gin-web/applications/gin-web/loader/resource"
	"github.com/mathiasXie/gin-web/applications/gin-web/middleware"
	"github.com/mathiasXie/gin-web/applications/tts-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/consts"
	"google.golang.org/grpc/metadata"
)

func InitRouter(ctx context.Context, r *gin.Engine) *gin.Engine {

	//r.Use(gin.Logger(), gin.Recovery())
	//engine.Use(middleware.Cors())

	userGroup := r.Group("/user")
	{
		userHandler := handler.UserHandler{
			Service: service.NewUserService(ctx),
		}
		userGroup.GET("list", func(c *gin.Context) {
			userHandler.GetUsers(c)
		})
		userGroup.POST("register", func(c *gin.Context) {
			userHandler.AddUser(c)
		})
		userGroup.POST("login", func(c *gin.Context) {
			//userHandler.Login(c)
		})
	}

	sysHandler := handler.SysHandler{
		Service: service.NewSysService(ctx),
	}
	managementGroup := r.Group("/v1/management")
	managementGroup.Use(middleware.AuthRoleMiddleware)
	{
		managementGroup.POST("routeArco", func(c *gin.Context) {
			sysHandler.RouteArco(c)
		})
		managementGroup.GET("route/list", func(c *gin.Context) {
			sysHandler.RouteList(c)
		})
		managementGroup.POST("route/update", func(c *gin.Context) {
			sysHandler.RouteUpdate(c)
		})
		managementGroup.POST("route/delete", func(c *gin.Context) {
			sysHandler.RouteDelete(c)
		})

		managementGroup.GET("user/list", func(c *gin.Context) {
			sysHandler.UserList(c)
		})
		managementGroup.POST("user/update", func(c *gin.Context) {
			sysHandler.UserUpdate(c)
		})
		managementGroup.POST("user/delete", func(c *gin.Context) {
			sysHandler.UserDelete(c)
		})
		managementGroup.POST("user/reset_password", func(c *gin.Context) {
			sysHandler.UserResetPassword(c)
		})
		managementGroup.GET("getUserInfo", func(c *gin.Context) {
			sysHandler.GetUserInfo(c)
		})

		managementGroup.POST("login", func(c *gin.Context) {
			sysHandler.Login(c)
		})
	}
	r.POST("/v1/login", func(ctx *gin.Context) {
		sysHandler.Login(ctx)
	})

	r.POST("/v1/verify", func(ctx *gin.Context) {
		verifyHandler := handler.VerifyHandler{}
		verifyHandler.MakeVerify(ctx)
	})

	r.POST("/sse", func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Content-Type", "text/event-stream")
		ctx.Writer.Header().Add("Cache-Control", "no-cache")
		ctx.Writer.Header().Add("Connection", "keep-alive")

		i := 0
		ctx.Stream(func(w io.Writer) bool {
			message := fmt.Sprintf("Server time: %s", time.Now().Format(time.DateTime))
			// 发送一块数据
			ctx.SSEvent("data", message)
			if i == 10 {
				return false
			}
			i++
			return true
		})
	})

	r.POST("/tts", func(ctx *gin.Context) {
		text := ctx.PostForm("text")

		//简单进行一下测试
		client := resource.GetResource().TTSRpcClient.TTSServiceClient

		// 创建带有元数据的上下文
		trace_id, _ := ctx.Value(consts.LogID).(string)
		md := metadata.Pairs("trace_id", trace_id)
		rpcCtx := metadata.NewOutgoingContext(ctx, md)

		resp, err := client.TextToSpeech(rpcCtx, &proto.TextToSpeechRequest{
			Provider: proto.Provider_MICROSOFT,
			VoiceId:  "zh-CN-XiaochenNeural",
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

	r.GET("/ws", func(ctx *gin.Context) {
		// 升级 HTTP 连接为 WebSocket 连接
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade connection:", err)
			return
		}
		defer conn.Close()

		for {
			// 读取客户端发送的消息
			messageType, p, err := conn.ReadMessage()

			if err != nil && err.Error() != "websocket: close 1000 (normal)" {
				log.Println("Error reading message:", err)
				break
			}
			// 打印接收到的消息
			log.Printf("Received message: %s", string(p))

			// 向客户端发送响应消息
			err = conn.WriteMessage(messageType, p)
			if err != nil {
				log.Println("Error writing message:", err)
				break
			}
		}
	})

	return r
}
