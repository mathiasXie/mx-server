package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAILLM struct {
	config Config
	client *openai.Client
}

func NewOpenAILLM(config Config) (LLMProvider, error) {

	if config.APIKey == "" {
		return nil, fmt.Errorf("openai llm requires an API key")
	}
	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL(config.BaseURL),
	)
	return &OpenAILLM{
		config: config,
		client: client,
	}, nil
}

// ChatStream 实现流式对话
func (a *OpenAILLM) ChatStream(ctx context.Context, modelID string, messages []*proto.ChatMessage, respChan chan<- ChatStreamResponse) error {
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
	m, _ := json.Marshal(chatMessages)
	logger.CtxInfo(ctx, "chat messages:", string(m))
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
	logger.CtxInfo(ctx, "OpenAILLM[ChatStream] 流式对话结束: ", totalContentStr)
	if stream.Err() != nil {
		return stream.Err()
	}

	// 发送结束标记
	respChan <- ChatStreamResponse{
		IsEnd: true,
	}
	return nil
}

// ChatNoStream 实现一次性返回对话
func (a *OpenAILLM) ChatNoStream(ctx context.Context, modelID string, messages []*proto.ChatMessage) (string, error) {
	// 构建消息列表
	startTime := time.Now()
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
	// 创建请求
	chatCompletion, err := a.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F(chatMessages),
			Model:    openai.F(modelID),
		},
	)
	if err != nil {
		return "", err
	}

	var resMessage string
	for _, choice := range chatCompletion.Choices {
		resMessage = fmt.Sprintf("%s%s", resMessage, choice.Message.Content)
	}
	logger.CtxInfo(ctx, "[ChatNoStream] 对话结束:", resMessage, ";耗时:", time.Since(startTime).Seconds())
	return resMessage, nil
}

func (a *OpenAILLM) GetModelList(ctx context.Context) ([]string, error) {
	return a.config.Models, nil
}
