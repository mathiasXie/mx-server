package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
)

func (h *ChatHandler) handlerHelloMessage(rpcCtx context.Context, chatRequest *dto.ChatRequest, conn *websocket.Conn) error {
	resp2, _ := json.Marshal(dto.ChatResponse{
		Type:      dto.ChatTypeHello,
		State:     dto.ChatStateStart,
		SessionID: "09fc016d-3425-462d-a92c-782ef53dacca",
		Text:      "你好，我是小智，一个智能助手。",
		Emotion:   "happy",
	})
	// 向客户端发送响应消息
	err := conn.WriteMessage(websocket.TextMessage, resp2)
	if err != nil {
		log.Println("Error writing message:", err)
		return err
	}
	return nil
}
