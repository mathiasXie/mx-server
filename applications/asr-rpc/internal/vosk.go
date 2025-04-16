package asr

import (
	"context"
	"fmt"
	"net/http"
)

// VoskASR 实现Vosk的ASR服务
type VoskASR struct {
	config Config
	client *http.Client
}

// NewVoskASR
func NewVoskASR(config Config) (ASRProvider, error) {

	if config.LibPath == "" {
		return nil, fmt.Errorf("vosk asr requires a lib path")
	}
	return &VoskASR{
		config: config,
		client: &http.Client{},
	}, nil
}

// SpeechToText 实现ASRProvider接口
func (v *VoskASR) SpeechToText(ctx context.Context, audioData []byte) (string, error) {

	return "你好啊", nil
}
