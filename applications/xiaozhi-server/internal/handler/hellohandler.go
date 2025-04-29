package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	llm_proto "github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/mgutz/ansi"
)

func (h *ChatHandler) handlerHelloMessage(chatRequest *dto.ChatRequest) error {

	resp2, _ := json.Marshal(dto.ChatResponse{
		Type:        dto.ChatTypeHello,
		State:       dto.ChatStateStart,
		SessionID:   h.sessionID,
		Emotion:     "happy",
		Transport:   "websocket",
		AudioParams: chatRequest.AudioParams,
	})
	log.Println("收到hello消息:", chatRequest)
	h.print("收到hello消息:", "green")
	// 向客户端发送响应消息
	err := h.conn.WriteMessage(websocket.TextMessage, resp2)
	if err != nil {
		log.Println("Error writing message:", err)
		return err
	}
	return nil
	// h.userInfo.Device = &dto.DeviceInfo{
	// 	DeviceId:   chatRequest.DeviceID,
	// 	DeviceName: chatRequest.DeviceName,
	// 	DeviceMac:  chatRequest.DeviceMac,
	// 	Token:      chatRequest.Token,
	// }
	// h.print(fmt.Sprintf("收到hello消息: %+v", chatRequest), "green")
	// // 检查设备是否存在
	// device, err := h.deviceService.GetDeviceByMac(chatRequest.DeviceMac)
	// if err != nil {
	// 	return err
	// }
	// if device == nil || device.RoleId == 0 {
	// 	h.print("设备不存在,前往绑定设备,会自动创建", "yellow")
	// 	h.deviceBindHandler()
	// } else if device.RoleId > 0 {
	// 	userInfo, err := h.userService.GetUserByRoleId(int(device.RoleId))
	// 	if err != nil {
	// 		h.print(fmt.Sprintf("获取用户信息失败:%s", err.Error()), "red")
	// 		return err
	// 	}
	// 	if userInfo != nil {
	// 		h.userInfo = userInfo
	// 		h.userInfo.Device = &dto.DeviceInfo{
	// 			Id:         int(device.Id),
	// 			DeviceId:   device.DeviceId,
	// 			DeviceName: device.DeviceName,
	// 			DeviceMac:  device.DeviceMac,
	// 			Token:      device.Token,
	// 		}
	// 		prompt := fmt.Sprintf("%s\n%s", config.Instance.Provider.PromptPrefix, h.userInfo.Role.RoleDesc)
	// 		h.generateChatContext(prompt)
	// 	}
	// }

	return nil
}

func (h *ChatHandler) checkDeviceBindStatus() bool {
	return h.userInfo.ID != 0
}
func (h *ChatHandler) print(text string, color string) {
	if config.Instance.RunMode != "prod" {
		printText := fmt.Sprintf("[%d] %s", h.userInfo.ID, text)
		log.Println(ansi.Color(printText, color))
	}
}

// 生成最近5条聊天记录
func (h *ChatHandler) generateChatContext(prompt string) {

	// h.userInfo.ChatMessages = append(h.userInfo.ChatMessages, &llm_proto.ChatMessage{
	// 	Role:    llm_proto.ChatMessageRole(llm_proto.ChatMessageRole_value[role]),
	// 	Content: text,
	// })

	llmMessages := make([]*llm_proto.ChatMessage, 1)
	llmMessages[0] = &llm_proto.ChatMessage{
		Role:    llm_proto.ChatMessageRole_SYSTEM,
		Content: prompt,
	}

	// 从数据库中获取最近5条聊天记录
	chatRecords, err := h.messageService.GetChatRecords(h.userInfo.ID, h.userInfo.Device.Id, 5)
	if err != nil {
		logger.CtxError(h.rpcCtx, "[ChatHandler]generateChatContext获取聊天记录失败:", err)
		return
	}

	for _, chatRecord := range chatRecords {
		llmMessages = append(llmMessages, &llm_proto.ChatMessage{
			Role:    llm_proto.ChatMessageRole(llm_proto.ChatMessageRole_value[chatRecord.Role]),
			Content: chatRecord.Message,
		})
	}

	h.userInfo.ChatMessages = llmMessages
	// for i, chatMessage := range h.userInfo.ChatMessages {
	// 	llmMessages[i+1] = &llm_proto.ChatMessage{
	// 		Role:    llm_proto.ChatMessageRole(llm_proto.ChatMessageRole_value[chatMessage.Role]),
	// 		Content: chatMessage.Content,
	// 	}
	// }
	// go func(role string, text string) {
	// 	err := h.messageService.StoreChatRecord(h.userInfo.ID, h.userInfo.Device.Id, role, text)
	// 	if err != nil {
	// 		logger.CtxError(h.rpcCtx, "[ChatHandler]storeChatRecord存储聊天记录失败:", err)
	// 	}
	// }(role, text)

	// return llmMessages
}
