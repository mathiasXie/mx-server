package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/mathiasXie/cloud_config"
	llm_proto "github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/mgutz/ansi"
)

func (h *ChatHandler) Authenticate(header http.Header, query func(key string) string) error {

	// 从请求参数中获取设备id
	deviceID := header.Get("device-id")
	if deviceID == "" {
		deviceID = query("device-id")
	}
	if deviceID == "" {
		log.Println(ansi.Color("无法从请求头或URL查询参数中获取device-id", "red"))
		return errors.New("无法从请求头或URL查询参数中获取device-id")
	}

	// 进行用户认证
	err := h.AuthUser(deviceID)
	if err != nil {
		log.Println(ansi.Color(fmt.Sprintf("用户认证失败:%v", err), "red"))
		return err
	}
	// 进行设备认证
	return nil
}

func (h *ChatHandler) AuthUser(deviceID string) error {
	device, err := h.deviceService.GetDeviceByMac(deviceID)
	if err != nil {
		return err
	}
	if device == nil || device.RoleId == 0 {
		return errors.New("设备不存在或未绑定用户")
	}

	user, err := h.userService.GetUserByRoleId(int(device.RoleId))
	if err != nil {
		return err
	}
	h.userInfo = user
	h.device = &dto.Device{
		Id:        int(device.Id),
		DeviceMac: device.DeviceMac,
		Language:  device.Language,
	}

	// 生成聊天上下文
	promptPrefix := cloud_config.GetConfig("provider")["prompt_prefix"]
	if promptPrefix == "" {
		promptPrefix = config.Instance.Provider.PromptPrefix
	}
	prompt := fmt.Sprintf("%s\n%s", promptPrefix, h.userInfo.Role.RoleDesc)
	h.generateChatContext(prompt)
	return nil
}

// 生成最近5条聊天记录
func (h *ChatHandler) generateChatContext(prompt string) {

	llmMessages := make([]*llm_proto.ChatMessage, 1)
	llmMessages[0] = &llm_proto.ChatMessage{
		Role:    llm_proto.ChatMessageRole_SYSTEM,
		Content: prompt,
	}

	// 从数据库中获取最近5条聊天记录
	chatRecords, err := h.messageService.GetChatRecords(h.userInfo.ID, h.device.Id, 5)
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
}
