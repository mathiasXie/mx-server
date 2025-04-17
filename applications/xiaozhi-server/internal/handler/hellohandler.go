package handler

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
)

func (h *ChatHandler) handlerHelloMessage(chatRequest *dto.ChatRequest) error {
	resp2, _ := json.Marshal(dto.ChatResponse{
		Type:      dto.ChatTypeHello,
		State:     dto.ChatStateStart,
		SessionID: h.sessionID,
		Text:      "你好，我是小智，一个智能助手。",
		Emotion:   "happy",
	})

	// 向客户端发送响应消息
	err := h.conn.WriteMessage(websocket.TextMessage, resp2)
	if err != nil {
		log.Println("Error writing message:", err)
		return err
	}
	return nil
}
