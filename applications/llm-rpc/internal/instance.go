package llm

import (
	"context"

	"github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

// 初始化各个LLM提供商
func NewInstanceLLMProvider() map[proto.LLMProvider]LLMProvider {
	llmProviders := map[proto.LLMProvider]LLMProvider{}

	if len(config.Instance.LLM.VolcEngine.APIKey) > 0 {
		volcEngineProvider, err := NewVolcEngineLLM(Config{
			APIKey:         config.Instance.LLM.VolcEngine.APIKey,
			BaseURL:        config.Instance.LLM.VolcEngine.BaseURL,
			Models:         config.Instance.LLM.VolcEngine.Models,
			DefaultModelID: config.Instance.LLM.VolcEngine.DefaultModelID,
		})
		if err != nil {
			logger.CtxError(context.Background(), "failed to create llm provider: %v", err)
			return nil
		}
		llmProviders[proto.LLMProvider_VOLCENGINE] = volcEngineProvider
	}
	if len(config.Instance.LLM.Aliyun.APIKey) > 0 {
		aliyunProvider, err := NewAliyunLLM(Config{
			APIKey:         config.Instance.LLM.Aliyun.APIKey,
			BaseURL:        config.Instance.LLM.Aliyun.BaseURL,
			DefaultModelID: config.Instance.LLM.Aliyun.DefaultModelID,
			Models:         config.Instance.LLM.Aliyun.Models,
		})
		if err != nil {
			logger.CtxError(context.Background(), "failed to create llm provider: %v", err)
			return nil
		}
		llmProviders[proto.LLMProvider_ALIYUN] = aliyunProvider
	}
	if len(config.Instance.LLM.GLM.APIKey) > 0 {
		glmProvider, err := NewGLMLLM(Config{
			APIKey:         config.Instance.LLM.GLM.APIKey,
			BaseURL:        config.Instance.LLM.GLM.BaseURL,
			DefaultModelID: config.Instance.LLM.GLM.DefaultModelID,
			Models:         config.Instance.LLM.GLM.Models,
		})
		if err != nil {
			logger.CtxError(context.Background(), "failed to create llm provider: %v", err)
			return nil
		}
		llmProviders[proto.LLMProvider_GLM] = glmProvider
	}
	return llmProviders
}
