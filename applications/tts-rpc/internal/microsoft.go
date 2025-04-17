package tts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

// MicrosoftTTS 实现微软的TTS服务
type MicrosoftTTS struct {
	config Config
	client *http.Client
}

type microVoices struct {
	VoiceID   string `json:"ShortName"`
	Language  string `json:"Locale"`
	VoiceName string `json:"DisplayName"`
}

// NewMicrosoftTTS 创建新的微软TTS实例
func NewMicrosoftTTS(config Config) (TTSProvider, error) {

	if config.APIKey == "" {
		return nil, fmt.Errorf("microsoft tts requires an API key")
	}
	return &MicrosoftTTS{
		config: config,
		client: &http.Client{},
	}, nil
}

// TextToSpeech 实现TTSProvider接口
func (m *MicrosoftTTS) TextToSpeech(ctx context.Context, text string, language string, voiceID string, speed float32, pitch float32) ([]byte, string, int32, int32, error) {

	/*
		<speak version='1.0' xml:lang='zh-CN'>
			<voice xml:lang='zh-CN' xml:gender='Female' name='zh-CN-XiaochenNeural'>
				我正在测试这个功能
			</voice>
		</speak>
	*/
	// 构建SSML
	ssml := fmt.Sprintf(`
		<speak version='1.0' xml:lang='%s'>
			<voice xml:lang='%s' xml:gender='Male' name='%s'>%s</voice>
		</speak>`, language, language, voiceID, text)

	// 构建请求
	req, err := http.NewRequestWithContext(ctx, "POST", m.config.Endpoint+"/cognitiveservices/v1", strings.NewReader(ssml))
	if err != nil {
		logger.CtxError(ctx, "failed to create request: %v", err)
		return nil, "", 0, 0, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", m.config.APIKey)
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("X-Microsoft-OutputFormat", "audio-16khz-128kbitrate-mono-mp3")

	// 发送请求
	resp, err := m.client.Do(req)
	if err != nil {
		logger.CtxError(ctx, "failed to send request: %v", err)
		return nil, "", 0, 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.CtxError(ctx, "request failed with status: %d", resp.StatusCode)
		return nil, "", 0, 0, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	// 读取响应
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.CtxError(ctx, "failed to read response: %v", err)
		return nil, "", 0, 0, fmt.Errorf("failed to read response: %v", err)
	}
	logger.CtxInfo(ctx, fmt.Sprintf("Audio data length: %d bytes", len(audioData)))

	return audioData, "mp3", 16000, 1, nil
}

func (m *MicrosoftTTS) TextToSpeechStream(ctx context.Context, text string, language string, voiceID string, respChan chan<- TTSStreamResponse) error {
	start := time.Now()
	region := "southeastasia"
	config, err := speech.NewSpeechConfigFromSubscription(m.config.APIKey, region)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return &json.MarshalerError{}
	}
	defer config.Close()

	config.SetSpeechSynthesisOutputFormat(common.Audio16Khz32KBitRateMonoMp3)

	speechSynthesizer, err := speech.NewSpeechSynthesizerFromConfig(config, nil)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return err
	}
	defer speechSynthesizer.Close()

	for {
		text = strings.TrimSuffix(text, "\n")
		if len(text) == 0 {
			break
		}

		task := speechSynthesizer.StartSpeakingTextAsync(text)
		var outcome speech.SpeechSynthesisOutcome
		select {
		case outcome = <-task:
		case <-time.After(60 * time.Second):
			fmt.Println("Timed out")
			return errors.New("Timed out")
		}
		defer outcome.Close()
		if outcome.Error != nil {
			fmt.Println("Got an error: ", outcome.Error)
			return outcome.Error
		}

		// in most case we want to streaming receive the audio to lower the latency,
		// we can use AudioDataStream to do so.
		stream, err := speech.NewAudioDataStreamFromSpeechSynthesisResult(outcome.Result)
		defer stream.Close()
		if err != nil {
			fmt.Println("Got an error: ", err)
			return err
		}

		var all_audio []byte
		audio_chunk := make([]byte, 10240)
		i := 1
		for {
			n, err := stream.Read(audio_chunk)

			if err == io.EOF {
				respChan <- TTSStreamResponse{
					IsEnd: true,
				}
				logger.CtxInfo(ctx, "[MicrosoftTTS]TextToSpeechStream输出结束:", text, ",耗时：", time.Since(start).Milliseconds())
				break
			}

			// 发送响应
			respChan <- TTSStreamResponse{
				Audio:  audio_chunk[:n],
				Sample: 16000,
				Format: "mp3",
			}
			all_audio = append(all_audio, audio_chunk[:n]...)
			i++
		}

		//fmt.Printf("Read [%d] bytes from audio data stream.\n", len(all_audio))
	}
	return nil
}

func (m *MicrosoftTTS) VoicesList(ctx context.Context) ([]Voices, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", m.config.Endpoint+"/voices/list", nil)
	if err != nil {
		logger.CtxError(ctx, "failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", m.config.APIKey)
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("X-Microsoft-OutputFormat", "audio-16khz-128kbitrate-mono-mp3")

	resp, err := m.client.Do(req)
	if err != nil {
		logger.CtxError(ctx, "failed to send request: %v", err)
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.CtxError(ctx, "failed to read response: %v", err)
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var microVoices []microVoices
	err = json.Unmarshal(body, &microVoices)
	if err != nil {
		logger.CtxError(ctx, "failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	voices := make([]Voices, len(microVoices))
	for i, v := range microVoices {
		voices[i] = Voices{
			VoiceID:   v.VoiceID,
			Language:  v.Language,
			VoiceName: v.VoiceName,
		}
	}
	logger.CtxInfo(ctx, "voices: %v", voices)
	return voices, nil
}
