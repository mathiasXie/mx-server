package asr

import (
	"context"
)

// ASRProvider 定义了语音转文本服务的接口
type ASRProvider interface {
	// SpeechToText 将语音文本转换为文本
	SpeechToText(ctx context.Context, audioData []byte) (string, error)
}

// Config 定义了TTS服务的配置
type Config struct {
	APIKey   string            // API密钥
	Endpoint string            // 服务端点
	Token    string            // 令牌
	APIID    string            // 应用ID
	Options  map[string]string // 其他配置选项
	LibPath  string            // 库路径
	Model    string            // 模型
}

// NewASRProvider 创建新的ASR提供商实例
type NewASRProvider func(config Config) (ASRProvider, error)
