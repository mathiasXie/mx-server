package tts

import (
	"context"
	"os"
	"testing"
)

func TestAliyunTTS_TextToSpeech(t *testing.T) {
	// 创建微软TTS实例
	apiKey := "sk-69734b46cb6a4a24a0f29ab84f7451ae" //os.Getenv("TTS_API_KEY")

	ttsProvider, err := NewAliyunTTS(Config{
		APIKey: apiKey,
	})
	if err != nil {
		t.Fatalf("Failed to create Aliyun TTS instance: %v", err)
	}

	// 测试文本转语音
	text := "测试文本转语音"
	language := "zh-CN"
	voiceID := "longxiaochun"
	speed := 1.0
	pitch := 0.0

	audioData, format, sampleRate, channels, err := ttsProvider.TextToSpeech(context.Background(), text, language, voiceID, float32(speed), float32(pitch))
	if err != nil {
		t.Fatalf("Failed to convert text to speech: %v", err)
	}
	// 将音频数据保存到文件
	filePath := "output22.mp3"
	err = os.WriteFile(filePath, audioData, 0644)
	if err != nil {
		t.Fatalf("Failed to save audio data to file: %v", err)
	}

	t.Logf("Audio data length: %d bytes", len(audioData))
	t.Logf("Audio format: %s", format)
	t.Logf("Sample rate: %d Hz", sampleRate)
	t.Logf("Channels: %d", channels)
}
