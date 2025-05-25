package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	asr_proto "github.com/mathiasXie/gin-web/applications/asr-rpc/proto/pb/proto"
	llm_proto "github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	tts_proto "github.com/mathiasXie/gin-web/applications/tts-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	internal_consts "github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/consts"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/service"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/loader"
	"github.com/mathiasXie/gin-web/consts"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/mathiasXie/gin-web/utils"
	"github.com/mgutz/ansi"
	"google.golang.org/grpc/metadata"
)

type ChatHandler struct {
	LLMClient        *llm_proto.LLMServiceClient
	TTSClient        *tts_proto.TTSServiceClient
	ASRClient        *asr_proto.ASRServiceClient
	asrAudio         []byte
	clientHaveVoice  bool
	clientVoiceStop  bool
	rpcCtx           context.Context
	conn             *websocket.Conn
	clientIP         string
	sessionID        string
	userInfo         *dto.UserInfo
	device           *dto.Device
	deviceService    *service.DeviceService
	messageService   *service.MessageService
	userService      *service.UserService
	roleService      *service.RoleService
	clientListenMode string
	clientAbort      bool

	vadSilenceThreshold  int64
	vadLastHaveVoiceTime int64
	vadAudioBuffer       []byte
}

func (h *ChatHandler) Chat(ctx *gin.Context, upgrader websocket.Upgrader) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.CtxError(ctx, "Failed to upgrade connection:", err)
		return
	}
	defer func() {
		log.Println(ansi.Color("服务端关闭连接", "red"))
		conn.Close()
	}()

	h.deviceService = service.NewDeviceService(ctx, loader.GetDB(ctx, true))
	h.messageService = service.NewMessageService(ctx, loader.GetDB(ctx, true))
	h.userService = service.NewUserService(ctx, loader.GetDB(ctx, true))
	h.roleService = service.NewRoleService(ctx, loader.GetDB(ctx, true))

	h.clientIP = utils.GetClientIP(ctx)

	// 进行设备认证
	err = h.Authenticate(ctx.Request.Header, ctx.Query)
	if err != nil {
		log.Println(ansi.Color("认证失败", "red"))
		return
	}
	// 认证通过,继续处理

	// 创建带有元数据的上下文
	trace_id, _ := ctx.Value(consts.LogID).(string)
	md := metadata.Pairs("trace_id", trace_id)
	rpcCtx := metadata.NewOutgoingContext(ctx, md)
	h.rpcCtx = rpcCtx

	h.conn = conn
	h.sessionID = trace_id

	h.clientListenMode = "manual"
	h.vadSilenceThreshold = 700
	h.resetVad()
	for {
		// 读取客户端发送的消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if err.Error() == "websocket: close 1000 (normal)" {
				h.print("客户端断开连接", "red")
			} else {
				h.print(fmt.Sprintf("读取消息失败: %s", err), "red")
				logger.CtxError(h.rpcCtx, "[ChatHandler]handlerMessage读取消息失败:", err)
			}
			break
		}

		h.handlerMessage(messageType, p)
	}
}

func (h *ChatHandler) handlerMessage(messageType int, p []byte) error {

	if messageType == websocket.TextMessage {
		//fmt.Printf("Received message: %s\n", string(p))

		chatRequest := dto.ChatRequest{}
		err := json.Unmarshal(p, &chatRequest)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			return err
		}
		//fmt.Printf("chatRequest: %+v\n", chatRequest)
		// 客户端连接问候消息
		if chatRequest.Type == dto.ChatTypeHello {
			return h.handlerHelloMessage(&chatRequest)
		} else if chatRequest.Type == dto.ChatTypeAbort {

		} else if chatRequest.Type == dto.ChatTypeListen {

			if chatRequest.Mode == "auto" {
				h.clientListenMode = "auto"
				//h.print("自动拾音模式", "blue")
			}
			// 收到音频开始消息
			if chatRequest.State == dto.ChatStateStart {
				//h.clientHaveVoice = true
				h.clientVoiceStop = false
			}
			// 收到音频结束消息
			if chatRequest.State == dto.ChatStateStop {
				//h.clientHaveVoice = true
				h.clientVoiceStop = true
				if len(h.asrAudio) > 0 {
					return h.handlerAudioMessage(nil)
				}
			}
			// 处理detect消息
			if chatRequest.State == dto.ChatStateDetect {
				return h.startToChat(chatRequest.Text)
			}
		} else if chatRequest.Type == dto.ChatTypeIOT {
		}
	} else if messageType == websocket.BinaryMessage {
		return h.handlerAudioMessage(p)
	}
	return nil
}

func (h *ChatHandler) startToChat(text string) error {

	h.print(fmt.Sprintf("用户说: %s", text), "blue")
	h.storeChatRecord("USER", text)

	// 进行意图判断
	indentText, needLLM, afterFun, err := h.IndentHandler()
	if err != nil {
		logger.CtxError(h.rpcCtx, "[ChatHandler]startToChat意图判断失败:", err)
		return err
	}
	full_text := ""
	if needLLM {

		if len(indentText) > 0 {
			h.userInfo.ChatMessages = append(h.userInfo.ChatMessages, &llm_proto.ChatMessage{
				Role:    llm_proto.ChatMessageRole(llm_proto.ChatMessageRole_value["SYSTEM"]),
				Content: indentText,
			})
		}

		// 大模型服务配置
		providerName, ok := llm_proto.LLMProvider_value[h.userInfo.Role.LLM]
		if !ok {
			logger.CtxError(h.rpcCtx, "[ChatHandler]startToChat大模型服务配置错误:", h.userInfo.Role.LLM)
			return internal_consts.RetLLMConfigError
		}
		provider := llm_proto.LLMProvider(providerName)

		// 调用LLM服务，获取回复
		resp, err := (*h.LLMClient).ChatStream(h.rpcCtx, &llm_proto.ChatRequest{
			Messages: h.userInfo.ChatMessages,
			Provider: provider,
			ModelId:  h.userInfo.Role.LLMModelId,
		})
		if err != nil {
			logger.CtxError(h.rpcCtx, "[ChatHandler]startToChat调用LLM服务失败:", err)
			return err
		}
		splitChars := "。？！；：?!"
		// 将字符集转换为 rune 切片

		processed_index := 0
		full_runes := make([]rune, 0)
		//服务器开始发送语音
		h.sendTTSMessage(dto.ChatStateStart)

		for {
			msg, err := resp.Recv()
			if err != nil {
				logger.CtxError(h.rpcCtx, "[ChatHandler]startToChat调用LLM服务失败:", err)
				return err
			}
			if msg.IsEnd {
				break
			}
			full_text = fmt.Sprintf("%s%s", full_text, msg.Content)

			// 将字符串转换为rune数组，以正确处理中文字符
			full_runes = append(full_runes, []rune(msg.Content)...)
			current_runes := full_runes[processed_index:] // 从未处理的位置开始，使用rune索引

			for i, p := range current_runes {
				if strings.Contains(splitChars, string(p)) {
					processed_index += i + 1

					respText := fmt.Sprintf("%s%s", string(current_runes[:i]), string(p))
					h.print(fmt.Sprintf("大模型说: %s", respText), "green")
					if len(respText) > 0 {
						// 向客户端发送音频消息
						h.sendAudioMessage(respText)
					}
				}
			}
		}
	} else {
		h.sendAudioMessage(indentText)
		full_text = indentText
	}
	//服务器语音传输结束
	//time.Sleep(3 * time.Second)
	h.print("服务器语音传输结束", "blue")
	h.sendTTSMessage(dto.ChatStateStop)
	h.storeChatRecord("ASSISTANT", full_text)
	if afterFun != nil {
		afterFun()
	}
	return nil
}

func (h *ChatHandler) sendTextMessage(text string, state dto.ChatState, chatType dto.ChatType) error {
	chatResponse := dto.ChatResponse{
		Type:      chatType,
		State:     state,
		SessionID: h.sessionID,
		Text:      text,
	}
	resp, _ := json.Marshal(chatResponse)
	err := h.conn.WriteMessage(websocket.TextMessage, resp)
	if err != nil {
		logger.CtxError(h.rpcCtx, "[ChatHandler]startToChat发送文本消息失败:", err)
		return err
	}
	return nil
}

// 存储聊天记录
func (h *ChatHandler) storeChatRecord(role string, text string) {

	// 确保会话中只存5条记录，否则删除第2条，要不然LLM输入token会太多
	if len(h.userInfo.ChatMessages) >= 5 {
		h.userInfo.ChatMessages = append(h.userInfo.ChatMessages[:1], h.userInfo.ChatMessages[2:]...)
	}

	h.userInfo.ChatMessages = append(h.userInfo.ChatMessages, &llm_proto.ChatMessage{
		Role:    llm_proto.ChatMessageRole(llm_proto.ChatMessageRole_value[role]),
		Content: text,
	})

	// 另开一个协程，存入DB
	go func() {
		err := h.messageService.StoreChatRecord(h.userInfo.ID, h.device.Id, role, text)
		if err != nil {
			logger.CtxError(h.rpcCtx, "[ChatHandler]storeChatRecord存储聊天记录失败:", err)
		}
	}()
}
