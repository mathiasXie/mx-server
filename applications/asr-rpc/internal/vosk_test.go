package asr

import (
	"context"
	"testing"
)

func TestVoskASR_SpeechToText(t *testing.T) {

	asrProvider, err := NewVoskASR(Config{
		Model: "../models/vosk-model-cn-0.22",
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
