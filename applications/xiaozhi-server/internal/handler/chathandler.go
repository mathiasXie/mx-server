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
	"github.com/mathiasXie/gin-web/consts"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"google.golang.org/grpc/metadata"
)

type ChatHandler struct {
	LLMClient       *llm_proto.LLMServiceClient
	TTSClient       *tts_proto.TTSServiceClient
	ASRClient       *asr_proto.ASRServiceClient
	asrAudio        []byte
	clientHaveVoice bool
	clientVoiceStop bool
	rpcCtx          context.Context
	conn            *websocket.Conn
	sessionID       string
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
	trace_id, _ := ctx.Value(consts.LogID).(string)
	md := metadata.Pairs("trace_id", trace_id)
	rpcCtx := metadata.NewOutgoingContext(ctx, md)
	h.rpcCtx = rpcCtx
	h.conn = conn
	h.sessionID = trace_id
	for {
		// 读取客户端发送的消息
		messageType, p, err := conn.ReadMessage()
		if err != nil && err.Error() != "websocket: close 1000 (normal)" {
			log.Println("Error reading message:", err)
			logger.CtxError(h.rpcCtx, "[ChatHandler]handlerMessage读取消息失败:", err)
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
		}
		// 收到音频开始消息
		if chatRequest.State == dto.ChatStateStart {
			h.clientHaveVoice = true
			h.clientVoiceStop = false
		}
		// 收到音频结束消息
		if chatRequest.State == dto.ChatStateStop {
			h.clientHaveVoice = true
			h.clientVoiceStop = true
			if len(h.asrAudio) > 0 {
				return h.handlerAudioMessage(nil)
			}
		}
		// 处理detect消息
		if chatRequest.State == dto.ChatStateDetect {
			return h.startToChat(chatRequest.Text)
		}
	} else if messageType == websocket.BinaryMessage {
		return h.handlerAudioMessage(p)
	}
	return nil
}

func (h *ChatHandler) startToChat(text string) error {

	log.Printf("\033[1;31m用户说: %s\033[0m\n", text)

	// 调用LLM服务，获取回复
	resp, err := (*h.LLMClient).ChatStream(h.rpcCtx, &llm_proto.ChatRequest{
		Messages: []*llm_proto.ChatMessage{
			{
				Role:    llm_proto.ChatMessageRole_SYSTEM,
				Content: "你是一个智能助手，请根据用户的问题给出回复,请简单明了，一句话搞定，不要使用复杂的句式",
			},
			{
				Role:    llm_proto.ChatMessageRole_USER,
				Content: text,
			},
		},
		Provider: llm_proto.LLMProvider_ALIYUN,
		ModelId:  "qwen-plus",
	})
	if err != nil {
		logger.CtxError(h.rpcCtx, "[ChatHandler]startToChat调用LLM服务失败:", err)
		return err
	}
	splitChars := "。？！；："
	// 将字符集转换为 rune 切片
	full_text := ""
	processed_index := 0
	full_runes := make([]rune, 0)
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
				log.Printf("\033[1;32m大模型说: %s\033[0m\n", respText)
				if len(respText) > 0 {
					// 向客户端发送开始文本消息
					err = h.sendTextMessage(respText, dto.ChatStateSentenceStart, dto.ChatTypeTTS)
					if err != nil {
						break
					}
					// 向客户端发送音频消息
					err = h.sendAudioMessage(respText)
					if err != nil {
						break
					}
					// 向客户端发送结束文本消息
					err = h.sendTextMessage(respText, dto.ChatStateSentenceEnd, dto.ChatTypeTTS)
					if err != nil {
						break
					}
				}
			}
		}
	}
	h.sendTTSMessage(dto.ChatStateStop)

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
