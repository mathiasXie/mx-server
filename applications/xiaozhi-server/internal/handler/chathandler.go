package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	asrAudio        []byte
	clientHaveVoice bool
	clientVoiceStop bool
	rpcCtx          context.Context
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
	for {
		// 读取客户端发送的消息
		messageType, p, err := conn.ReadMessage()
		if err != nil && err.Error() != "websocket: close 1000 (normal)" {
			log.Println("Error reading message:", err)
			break
		}
		// 打印接收到的消息
		fmt.Println("--------------------------------")
		fmt.Println("消息类型", messageType)
		fmt.Println("--------------------------------")
		h.handlerMessage(rpcCtx, messageType, p, conn)

	}
}

func (h *ChatHandler) handlerMessage(rpcCtx context.Context, messageType int, p []byte, conn *websocket.Conn) error {

	if messageType == websocket.TextMessage {
		fmt.Printf("Received message: %s\n", string(p))

		chatRequest := dto.ChatRequest{}
		err := json.Unmarshal(p, &chatRequest)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			return err
		}
		// 客户端连接问候消息
		if chatRequest.Type == dto.ChatTypeHello {
			return h.handlerHelloMessage(rpcCtx, &chatRequest, conn)
		}
		// 收到音频开始消息
		if chatRequest.State == dto.ChatStateStart {
			h.clientHaveVoice = true
			h.clientVoiceStop = false
			fmt.Println("--------------------------------")
			fmt.Println("收到音频开始消息")
			fmt.Println("--------------------------------")
		}
		// 收到音频结束消息
		if chatRequest.State == dto.ChatStateStop {
			h.clientHaveVoice = true
			h.clientVoiceStop = true
			if len(h.asrAudio) > 0 {
				return h.handlerAudioMessage(rpcCtx, nil, conn)
			}
		}
		// 处理detect消息
		if chatRequest.State == dto.ChatStateDetect {
			return h.handlerTextMessage(rpcCtx, &chatRequest, conn)
		}
	} else if messageType == websocket.BinaryMessage {

		return h.handlerAudioMessage(rpcCtx, p, conn)
	}
	return nil
}
