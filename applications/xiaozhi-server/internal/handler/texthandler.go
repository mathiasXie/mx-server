package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	tts_proto "github.com/mathiasXie/gin-web/applications/tts-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	audio_utils "github.com/mathiasXie/gin-web/utils/audio"
)

func (h *ChatHandler) handlerTextMessage(rpcCtx context.Context, chatRequest *dto.ChatRequest, conn *websocket.Conn) error {

	if chatRequest.Type == dto.ChatTypeHello {

	} else {

		for i := 1; i <= 1; i++ {
			text := "你好啊，祝你每天好心情"
			chatResponse := dto.ChatResponse{
				Type:      dto.ChatTypeTTS,
				State:     dto.ChatStateSentenceStart,
				SessionID: "09fc016d-3425-462d-a92c-782ef53dacca",
				Text:      text,
				Emotion:   "happy",
			}
			resp, _ := json.Marshal(chatResponse)
			err := conn.WriteMessage(websocket.TextMessage, resp)
			if err != nil {
				log.Println("Error writing message:", err)
				return err
			}

			//发送一段语音给客户端
			ttsResp, err := (*h.TTSClient).TextToSpeechStream(rpcCtx, &tts_proto.TextToSpeechRequest{
				Provider: tts_proto.Provider_MICROSOFT,
				VoiceId:  "zh-CN-Xiaochen:DragonHDLatestNeural",
				Language: "zh-CN",
				Text:     text,
			})
			if err != nil {
				log.Println("Error writing message:", err)
				break
			} else {
				for {
					msg, err := ttsResp.Recv()
					if err != nil {
						log.Println("Error writing message:", err)
						break
					}
					if len(msg.AudioData) > 0 {
						tempAudioFile, err := os.CreateTemp("", "to_tts_*.mp3")
						if err != nil {
							fmt.Println("无法创建临时音频文件: ", err)
							break
						}
						err = os.WriteFile(tempAudioFile.Name(), msg.AudioData, 0644)
						if err != nil {
							log.Println("Error writing message:", err)
							break
						}

						opusData, _, err := audio_utils.AudioToOpusData(tempAudioFile.Name(), 16000, 1)
						if err != nil {
							log.Println("Error writing message:", err)
							break
						}
						os.Remove(tempAudioFile.Name())
						sendAudioData(opusData, conn)
					}

					if msg.IsEnd {
						chatResponse := dto.ChatResponse{
							Type:      dto.ChatTypeTTS,
							State:     dto.ChatStateSentenceEnd,
							SessionID: "09fc016d-3425-462d-a92c-782ef53dacca",
						}
						resp, _ := json.Marshal(chatResponse)
						err := conn.WriteMessage(websocket.TextMessage, resp)
						if err != nil {
							log.Println("Error writing message:", err)
							return err
						}
						break
					}
				}
			}

			log.Println("success")
			//}
		}
	}
	return nil
}

func sendAudioData(opusData []audio_utils.AudioByte, conn *websocket.Conn) {
	for _, data := range opusData {
		err := conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Println("Error writing message:", err)
			break
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
