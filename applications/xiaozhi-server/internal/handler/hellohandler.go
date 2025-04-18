package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mgutz/ansi"
)

func (h *ChatHandler) handlerHelloMessage(chatRequest *dto.ChatRequest) error {

	resp2, _ := json.Marshal(dto.ChatResponse{
		Type:      dto.ChatTypeHello,
		State:     dto.ChatStateStart,
		SessionID: h.sessionID,
		Text:      "你好，我是小智，一个智能助手。",
		Emotion:   "happy",
	})
	h.userInfo.Mac = chatRequest.DeviceMac
	h.print("收到hello消息:", "green")
	// 向客户端发送响应消息
	err := h.conn.WriteMessage(websocket.TextMessage, resp2)
	if err != nil {
		log.Println("Error writing message:", err)
		return err
	}
	return nil
}

func (h *ChatHandler) print(text string, color string) {
	if config.Instance.RunMode != "prod" {
		printText := fmt.Sprintf("[%s] %s", h.userInfo.Mac, text)
		log.Println(ansi.Color(printText, color))
	}
}
