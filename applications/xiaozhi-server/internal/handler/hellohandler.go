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

	h.userInfo.Device = &dto.DeviceInfo{
		DeviceId:   chatRequest.DeviceID,
		DeviceName: chatRequest.DeviceName,
		DeviceMac:  chatRequest.DeviceMac,
		Token:      chatRequest.Token,
	}

	// 检查设备是否存在
	device, err := h.deviceService.GetDeviceByMac(chatRequest.DeviceMac)
	if err != nil {
		return err
	}
	fmt.Println("device", device)
	if device == nil || device.RoleId == 0 {
		h.print("设备不存在,前往绑定设备,会自动创建", "yellow")
		return h.deviceBindHandler()
	} else if device.RoleId > 0 {
		userInfo, err := h.userService.GetUserByRoleId(int(device.RoleId))
		if err != nil {
			h.print(fmt.Sprintf("获取用户信息失败:%s", err.Error()), "red")
			return err
		}
		if userInfo != nil {
			h.userInfo = userInfo
		}
	}

	resp2, _ := json.Marshal(dto.ChatResponse{
		Type:      dto.ChatTypeHello,
		State:     dto.ChatStateStart,
		SessionID: h.sessionID,
		Text:      "你好，我是小智，一个智能助手。",
		Emotion:   "happy",
	})
	log.Println("收到hello消息:", chatRequest)
	h.print("收到hello消息:", "green")
	// 向客户端发送响应消息
	err = h.conn.WriteMessage(websocket.TextMessage, resp2)
	if err != nil {
		log.Println("Error writing message:", err)
		return err
	}
	return nil
}

func (h *ChatHandler) deviceBindHandler() error {
	bindCode := h.deviceService.GenerateBindCode(h.userInfo.Device)

	text := fmt.Sprintf("请前往用户中心绑定设备,验证码是%d", bindCode)
	h.sendAudioMessage(text)
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
