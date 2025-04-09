package router

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	function_proto "github.com/mathiasXie/gin-web/applications/function-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/handler"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/service"
	"github.com/mathiasXie/gin-web/applications/gin-web/loader/resource"
	"github.com/mathiasXie/gin-web/applications/gin-web/middleware"
	llm_proto "github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
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

	r.POST("/weather", func(ctx *gin.Context) {
		location := ctx.PostForm("location")

		//简单进行一下测试
		client := resource.GetResource().FunctionRpcClient.FunctionServiceClient
		resp, err := client.GetWeatherReport(ctx, &function_proto.GetWeatherReportRequest{
			Lang:     "zh",
			Location: location,
		})
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": err.Error(),
			})
		}
		ctx.JSON(200, gin.H{
			"message": "success",
			"data":    resp,
		})
	})

	r.POST("/llm", func(ctx *gin.Context) {
		request_content := ctx.PostForm("request_content")

		f_client := resource.GetResource().FunctionRpcClient.FunctionServiceClient
		f_resp, _ := f_client.GetWeatherReport(ctx, &function_proto.GetWeatherReportRequest{
			Lang:     "zh",
			Location: "上海",
		})

		client := resource.GetResource().LLMRpcClient.LLMServiceClient
		// 创建带有元数据的上下文
		trace_id, _ := ctx.Value(consts.LogID).(string)
		md := metadata.Pairs("trace_id", trace_id)
		rpcCtx := metadata.NewOutgoingContext(ctx, md)
		resp, err := client.ChatStream(rpcCtx, &llm_proto.ChatRequest{
			Messages: []*llm_proto.ChatMessage{
				{
					Role:    llm_proto.ChatMessageRole_SYSTEM,
					Content: f_resp.Report,
				},
				{
					Role:    llm_proto.ChatMessageRole_USER,
					Content: request_content,
				},
			},
			Provider: llm_proto.LLMProvider_ALIYUN,
			ModelId:  "qwen-plus",
		})
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": err.Error(),
			})
		}
		splitChars := "。？！；："
		// 将字符集转换为 rune 切片
		//splitRunes := string([]rune(splitChars))
		full_text := ""
		processed_index := 0
		full_runes := make([]rune, 0)
		for {
			msg, err := resp.Recv()
			if err != nil {
				ctx.JSON(500, gin.H{
					"message": err.Error(),
				})
			}
			//流式返回
			if msg.IsEnd {
				break
			}
			full_text = fmt.Sprintf("%s%s", full_text, msg.Content)
			fmt.Println("full_text", full_text)

			// 将字符串转换为rune数组，以正确处理中文字符
			full_runes = append(full_runes, []rune(msg.Content)...)
			current_runes := full_runes[processed_index:] // 从未处理的位置开始，使用rune索引
			fmt.Println("current_runes", string(current_runes))

			for i, p := range current_runes {
				if strings.Contains(splitChars, string(p)) {
					fmt.Printf("找到句号index:%d,processed_index:%d\n", i, processed_index)
					ctx.SSEvent("data", string(current_runes[:i]))
					processed_index += i + 1
				}
			}
			//processed_chars += len(current_text)
			// if strings.Contains(current_text, p) {
			// 	pos := strings.Index(current_text, p)
			// 	// 计算实际的rune数量，而不是字节数
			// 	posRunes := len([]rune(current_text[:pos+1])) // 包含标点符号
			// 	processed_chars += posRunes
			// 	ctx.SSEvent("data", current_text[:pos+1]) // 包含标点符号发送
			// 	break
			// }

		}
		fmt.Println("over_full_text", full_text)
	})

	r.GET("/ws_llm", func(ctx *gin.Context) {
		// 升级 HTTP 连接为 WebSocket 连接
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade connection:", err)
			return
		}
		defer conn.Close()

		f_client := resource.GetResource().FunctionRpcClient.FunctionServiceClient
		f_resp, _ := f_client.GetWeatherReport(ctx, &function_proto.GetWeatherReportRequest{
			Lang:     "zh",
			Location: "上海",
		})

		client := resource.GetResource().LLMRpcClient.LLMServiceClient
		// 创建带有元数据的上下文
		trace_id, _ := ctx.Value(consts.LogID).(string)
		md := metadata.Pairs("trace_id", trace_id)
		rpcCtx := metadata.NewOutgoingContext(ctx, md)

		for {
			// 读取客户端发送的消息
			messageType, p, err := conn.ReadMessage()
			if err != nil && err.Error() != "websocket: close 1000 (normal)" {
				log.Println("Error reading message:", err)
				break
			}

			resp, err := client.ChatStream(rpcCtx, &llm_proto.ChatRequest{
				Messages: []*llm_proto.ChatMessage{
					{
						Role:    llm_proto.ChatMessageRole_SYSTEM,
						Content: f_resp.Report,
					},
					{
						Role:    llm_proto.ChatMessageRole_USER,
						Content: string(p),
					},
				},
				Provider: llm_proto.LLMProvider_ALIYUN,
				ModelId:  "qwen-plus",
			})
			if err != nil {
				ctx.JSON(500, gin.H{
					"message": err.Error(),
				})
			}
			splitChars := "。？！；："
			// 将字符集转换为 rune 切片
			//splitRunes := string([]rune(splitChars))
			full_text := ""
			processed_index := 0
			full_runes := make([]rune, 0)
			for {
				msg, err := resp.Recv()
				if err != nil {
					ctx.JSON(500, gin.H{
						"message": err.Error(),
					})
				}
				//流式返回
				if msg.IsEnd {
					break
				}
				full_text = fmt.Sprintf("%s%s", full_text, msg.Content)

				// 将字符串转换为rune数组，以正确处理中文字符
				full_runes = append(full_runes, []rune(msg.Content)...)
				current_runes := full_runes[processed_index:] // 从未处理的位置开始，使用rune索引

				for i, p := range current_runes {
					if strings.Contains(splitChars, string(p)) {
						fmt.Printf("找到句号index:%d,processed_index:%d\n", i, processed_index)
						//ctx.SSEvent("data", string(current_runes[:i]))
						// 向客户端发送响应消息
						err = conn.WriteMessage(messageType, []byte(string(current_runes[:i])))
						if err != nil {
							log.Println("Error writing message:", err)
							break
						}
						processed_index += i + 1
					}
				}

			}

		}
	})
	return r
}
