package asr

import (
	"context"
	"testing"
)

func TestAliyunASR_SpeechToText(t *testing.T) {
	// 创建微软TTS实例
	apiKey := "sk-69734b46cb6a4a24a0f29ab84f7451ae" //os.Getenv("TTS_API_KEY")

	asrProvider, err := NewAliyunASR(Config{
		APIKey: apiKey,
	})
	if err != nil {
		t.Fatalf("Failed to create Aliyun TTS instance: %v", err)
	}

	text, err := asrProvider.SpeechToText(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to convert text to speech: %v", err)
	}

	t.Logf("语音识别结果: %s", text)
}
