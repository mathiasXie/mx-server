package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	llm_proto "github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	tts_proto "github.com/mathiasXie/gin-web/applications/tts-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

type ChatHandler struct {
	LLMClient llm_proto.LLMServiceClient
	TTSClient tts_proto.TTSServiceClient
}

func (h *ChatHandler) Chat(ctx *gin.Context, upgrader websocket.Upgrader) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.CtxError(ctx, "Failed to upgrade connection:", err)
		return
	}
	defer func() {
		fmt.Println("客户端断开连接")
		conn.Close()
	}()
	// client := resource.GetResource().LLMRpcClient.LLMServiceClient
	// 创建带有元数据的上下文
	// trace_id, _ := ctx.Value(consts.LogID).(string)
	// md := metadata.Pairs("trace_id", trace_id)
	// rpcCtx := metadata.NewOutgoingContext(ctx, md)

	for {
		// 读取客户端发送的消息
		messageType, p, err := conn.ReadMessage()
		if err != nil && err.Error() != "websocket: close 1000 (normal)" {
			log.Println("Error reading message:", err)
			break
		}
		// 打印接收到的消息
		log.Printf("Received message: %s", string(p))
		chatRequest := dto.ChatRequest{}
		err = json.Unmarshal(p, &chatRequest)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			break
		}
		if chatRequest.Type == dto.ChatTypeHello {
			resp2, _ := json.Marshal(dto.ChatResponse{
				Type:      dto.ChatTypeHello,
				State:     dto.ChatStateStart,
				SessionID: "09fc016d-3425-462d-a92c-782ef53dacca",
				Text:      "你好，我是小智，一个智能助手。",
				Emotion:   "happy",
			})
			// 向客户端发送响应消息
			err = conn.WriteMessage(messageType, resp2)
			if err != nil {
				log.Println("Error writing message:", err)
				break
			}
		} else {

			for i := 0; i < 1; i++ {
				chatResponse := dto.ChatResponse{
					Type:      dto.ChatTypeTTS,
					State:     dto.ChatStateSentenceStart,
					SessionID: "09fc016d-3425-462d-a92c-782ef53dacca",
					Text:      fmt.Sprintf("%d. 将鸡蛋打入碗中，加入适量盐，搅拌均匀", i),
					Emotion:   "happy",
				}
				resp, _ := json.Marshal(chatResponse)

				// 向客户端发送响应消息
				err = conn.WriteMessage(messageType, resp)
				if err != nil {
					log.Println("Error writing message:", err)
					break
				}

				//发送一段语音给客户端
				// ttsResp, err := h.TTSClient.TextToSpeech(rpcCtx, &tts_proto.TextToSpeechRequest{
				// 	Provider: tts_proto.Provider_ALIYUN,
				// 	VoiceId:  "longxiaochun",
				// 	Language: "zh-CN",
				// 	Text:     fmt.Sprintf("%d. 将鸡蛋打入碗中，加入适量盐，搅拌均匀", i),
				// })
				// if err != nil {
				// 	log.Println("Error writing message:", err)
				// 	break
				// } else {

				filePath := "./tts2.opus"
				file, err := os.Open(filePath)
				if err != nil {
					fmt.Println("打开文件时出错:", err)
					return
				}
				// 确保在函数结束时关闭文件
				defer file.Close()
				// 获取文件的信息，用于确定文件大小
				fileInfo, err := file.Stat()
				if err != nil {
					fmt.Println("获取文件信息时出错:", err)
					return
				}

				// 创建一个与文件大小相同的字节切片
				AudioData := make([]byte, fileInfo.Size())

				// 从文件中读取内容到 data 切片
				_, err = io.ReadFull(file, AudioData)
				if err != nil {
					if err == io.EOF {
						fmt.Println("文件读取不完整，可能文件损坏")
					} else {
						fmt.Println("读取文件时出错:", err)
					}
					return
				}

				// 向客户端发送响应消息
				// opusDatas, err := utils.AudioToOpusData(AudioData)
				// if err != nil {
				// 	fmt.Println("读取文件时出错:", err)
				// }
				err = conn.WriteMessage(websocket.BinaryMessage, AudioData)
				if err != nil {
					log.Println("Error writing message:", err)
					break
				}
				log.Println("success")
				//}
			}
		}

	}
}

/*
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
*/
