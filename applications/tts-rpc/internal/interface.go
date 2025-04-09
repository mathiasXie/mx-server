package tts

import (
	"context"
)

// TTSProvider 定义了文本转语音服务的接口
type TTSProvider interface {
	// TextToSpeech 将文本转换为语音
	TextToSpeech(ctx context.Context, text string, language string, voiceID string, speed float32, pitch float32) ([]byte, string, int32, int32, error)

	VoicesList(ctx context.Context) ([]Voices, error)
}

type Voices struct {
	VoiceID   string `json:"voice_id"`
	Language  string `json:"language"`
	VoiceName string `json:"voice_name"`
}

// Config 定义了TTS服务的配置
type Config struct {
	APIKey   string            // API密钥
	Endpoint string            // 服务端点
	Token    string            // 令牌
	APIID    string            // 应用ID
	Options  map[string]string // 其他配置选项
}

// NewTTSProvider 创建新的TTS提供商实例
type NewTTSProvider func(config Config) (TTSProvider, error)
