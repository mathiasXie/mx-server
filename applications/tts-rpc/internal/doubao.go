package tts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DoubaoTTS 实现豆包的TTS服务
type DoubaoTTS struct {
	config Config
	client *http.Client
}

// DoubaoRequest 豆包API请求结构
type DoubaoRequest struct {
	Text     string  `json:"text"`
	Language string  `json:"language"`
	VoiceID  string  `json:"voice_id"`
	Speed    float32 `json:"speed"`
	Pitch    float32 `json:"pitch"`
}

// NewDoubaoTTS 创建新的豆包TTS实例
func NewDoubaoTTS(config Config) (TTSProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("doubao tts requires an API key")
	}
	if config.Endpoint == "" {
		config.Endpoint = "https://api.doubao.com/v1/tts"
	}
	return &DoubaoTTS{
		config: config,
		client: &http.Client{},
	}, nil
}

// TextToSpeech 实现TTSProvider接口
func (d *DoubaoTTS) TextToSpeech(ctx context.Context, text string, language string, voiceID string, speed float32, pitch float32) ([]byte, string, int32, int32, error) {
	// 构建请求体
	reqBody := DoubaoRequest{
		Text:     text,
		Language: language,
		VoiceID:  voiceID,
		Speed:    speed,
		Pitch:    pitch,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, "", 0, 0, fmt.Errorf("failed to marshal request: %v", err)
	}

	// 构建请求
	req, err := http.NewRequestWithContext(ctx, "POST", d.config.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, "", 0, 0, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.config.APIKey))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, "", 0, 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", 0, 0, fmt.Errorf("request failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	// 读取响应
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", 0, 0, fmt.Errorf("failed to read response: %v", err)
	}

	return audioData, "mp3", 16000, 1, nil
}

func (d *DoubaoTTS) VoicesList(ctx context.Context) ([]Voices, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", d.config.Endpoint+"/voices/list", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.config.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var voices []Voices
	err = json.Unmarshal(body, &voices)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return voices, nil
}
