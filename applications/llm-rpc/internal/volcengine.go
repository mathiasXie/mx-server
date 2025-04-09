package llm

import (
	"context"

	"github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
)

type VolcEngineLLM struct {
	config Config
}

// NewVolcEngineLLM 创建新的火山引擎LLM实例
func NewVolcEngineLLM(config Config) (LLMProvider, error) {

	return nil, nil
	// if config.APIKey == "" {
	// 	return nil, fmt.Errorf("aliyun llm requires an API key")
	// }
	// return &VolcEngineLLM{
	// 	config: config,
	// }, nil
}

func (v *VolcEngineLLM) ChatStream(ctx context.Context, modelID string, messages []*proto.ChatMessage, respChan chan<- ChatStreamResponse) error {
	return nil
}

func (v *VolcEngineLLM) GetModelList(ctx context.Context) ([]string, error) {
	return v.config.Models, nil
}
