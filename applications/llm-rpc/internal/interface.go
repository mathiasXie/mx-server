package llm

import (
	"context"

	"github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
)

// LLMProvider 大语言模型服务接口
type LLMProvider interface {
	ChatStream(ctx context.Context, modelID string, messages []*proto.ChatMessage, respChan chan<- ChatStreamResponse) error
	GetModelList(ctx context.Context) ([]string, error)
}

type ChatStreamResponse struct {
	ID      string
	Content string
	IsEnd   bool
	Created int64
	Model   string
}

// Config 定义了LLM服务的配置
type Config struct {
	APIKey         string   // API密钥
	BaseURL        string   // 基础URL
	DefaultModelID string   // 默认模型ID
	Models         []string // 模型列表
}

// NewLLMProvider 创建新的LLM提供商实例
type NewLLMProvider func(config Config) (LLMProvider, error)
