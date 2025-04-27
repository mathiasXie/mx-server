package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type AliyunLLM struct {
	config Config
	client *openai.Client
}

// NewAliyunLLM 创建新的阿里云LLM实例
func NewAliyunLLM(config Config) (LLMProvider, error) {

	if config.APIKey == "" {
		return nil, fmt.Errorf("aliyun llm requires an API key")
	}
	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL(config.BaseURL),
	)
	return &AliyunLLM{
		config: config,
		client: client,
	}, nil
}

// ChatStream 实现流式对话
func (a *AliyunLLM) ChatStream(ctx context.Context, modelID string, messages []*proto.ChatMessage, respChan chan<- ChatStreamResponse) error {
	// 构建消息列表
	chatMessages := make([]openai.ChatCompletionMessageParamUnion, len(messages))
	for i, msg := range messages {
		switch msg.Role {
		case proto.ChatMessageRole_USER:
			chatMessages[i] = openai.UserMessage(msg.Content)
		case proto.ChatMessageRole_ASSISTANT:
			chatMessages[i] = openai.AssistantMessage(msg.Content)
		case proto.ChatMessageRole_SYSTEM:
			chatMessages[i] = openai.SystemMessage(msg.Content)
		}
	}
	logger.CtxInfo(ctx, "chat messages: %v", chatMessages)
	// 创建流式请求
	stream := a.client.Chat.Completions.NewStreaming(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F(chatMessages),
			Model:    openai.F(modelID),
		},
	)

	// 处理流式响应
	totalContent := []string{}
	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			totalContent = append(totalContent, chunk.Choices[0].Delta.Content)
			// 发送响应
			respChan <- ChatStreamResponse{
				ID:      chunk.ID,
				Content: chunk.Choices[0].Delta.Content,
				IsEnd:   chunk.Choices[0].FinishReason == openai.ChatCompletionChunkChoicesFinishReasonStop,
				Created: chunk.Created,
				Model:   chunk.Model,
			}
		}

	}
	totalContentStr := strings.Join(totalContent, "")
	logger.CtxInfo(ctx, "AliyunLLM[ChatStream] 流式对话结束: ", totalContentStr)
	if stream.Err() != nil {
		return stream.Err()
	}

	// 发送结束标记
	respChan <- ChatStreamResponse{
		IsEnd: true,
	}
	return nil
}

func (a *AliyunLLM) GetModelList(ctx context.Context) ([]string, error) {
	return a.config.Models, nil
}
