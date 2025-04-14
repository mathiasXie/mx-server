package tts

import (
	"context"
	"os"
	"testing"
)

func TestMicrosoftTTS_TextToSpeech(t *testing.T) {
	// 创建微软TTS实例
	apiKey := "8aKNEoi9njCkkJSZIKX4mHSXDCBl8HEtRlrP4ev4Zuw7X4J93CL2JQQJ99BCACqBBLyXJ3w3AAAYACOGJKg4" //os.Getenv("TTS_API_KEY")
	endpoint := "https://southeastasia.tts.speech.microsoft.com"

	ttsProvider, err := NewMicrosoftTTS(Config{
		APIKey:   apiKey,
		Endpoint: endpoint,
	})
	if err != nil {
		t.Fatalf("Failed to create Microsoft TTS instance: %v", err)
	}

	// 测试文本转语音
	text := "辛苦针对移动的NLU文档检查下，是不是每个意图都赋予了tag呀,另外，对方希望在NLU schema这个表，加一列tag字段"
	language := "zh-CN"
	voiceID := "zh-CN-XiaochenNeural"
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

func TestMicrosoftTTS_TextToSpeechStream(t *testing.T) {
	// 创建微软TTS实例
	apiKey := "8aKNEoi9njCkkJSZIKX4mHSXDCBl8HEtRlrP4ev4Zuw7X4J93CL2JQQJ99BCACqBBLyXJ3w3AAAYACOGJKg4" //os.Getenv("TTS_API_KEY")
	endpoint := "https://southeastasia.tts.speech.microsoft.com"

	ttsProvider, err := NewMicrosoftTTS(Config{
		APIKey:   apiKey,
		Endpoint: endpoint,
	})
	if err != nil {
		t.Fatalf("Failed to create Microsoft TTS instance: %v", err)
	}

	// 测试文本转语音
	text := "辛苦针对移动的NLU文档检查下，是不是每个意图都赋予了tag呀,另外，对方希望在NLU schema这个表，加一列tag字段"
	language := "zh-CN"
	voiceID := "zh-CN-XiaochenNeural"

	respChan := make(chan TTSStreamResponse)

	err = ttsProvider.TextToSpeechStream(context.Background(), text, language, voiceID, respChan)
	if err != nil {
		t.Fatalf("Failed to convert text to speech: %v", err)
	}
	// 处理流式响应
	for {
		select {
		case resp := <-respChan:
			// 将音频数据保存到文件
			err = os.WriteFile("output_stream.mp3", resp.Audio, 0644)
			if err != nil {
				t.Fatalf("Failed to save audio data to file: %v", err)
			}
		}
	}

}

func TestMicrosoftTTS_VoicesList(t *testing.T) {
	apiKey := "8aKNEoi9njCkkJSZIKX4mHSXDCBl8HEtRlrP4ev4Zuw7X4J93CL2JQQJ99BCACqBBLyXJ3w3AAAYACOGJKg4" //os.Getenv("TTS_API_KEY")
	endpoint := "https://southeastasia.tts.speech.microsoft.com/cognitiveservices"

	ttsProvider, err := NewMicrosoftTTS(Config{
		APIKey:   apiKey,
		Endpoint: endpoint,
	})
	if err != nil {
		t.Fatalf("Failed to create Microsoft TTS instance: %v", err)
	}

	voices, err := ttsProvider.VoicesList(context.Background())
	if err != nil {
		t.Fatalf("Failed to get voices list: %v", err)
	}

	t.Logf("Voices list: %v", voices)
}
