package asr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	vosk "github.com/alphacep/vosk-api/go"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

// VoskASR 实现Vosk的ASR服务
type VoskASR struct {
	config Config
	model  *vosk.VoskModel
	client *http.Client
}

// NewVoskASR
func NewVoskASR(config Config) (ASRProvider, error) {

	if config.Model == "" {
		return nil, fmt.Errorf("vosk asr requires a model path")
	}

	// 获取当前文件的调用栈信息
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get caller information")
	}

	// 获取文件所在的目录
	dir := filepath.Dir(filename)
	model, err := vosk.NewModel(fmt.Sprintf("%s/../models/%s", dir, config.Model))
	if err != nil {
		logger.CtxError(context.Background(), "[VoskASR]NewModel失败：", err)
		return nil, err
	}
	return &VoskASR{
		config: config,
		model:  model,
		client: &http.Client{},
	}, nil
}

// SpeechToText 实现ASRProvider接口
func (v *VoskASR) SpeechToText(ctx context.Context, audioPcmData []byte) (string, error) {
	startTime := time.Now()

	// we can check if word is in the vocabulary
	// fmt.Println(model.FindWord("air"))

	sampleRate := 16000.0
	rec, err := vosk.NewRecognizer(v.model, sampleRate)
	if err != nil {
		logger.CtxError(ctx, "[VoskASR]NewRecognizer失败：", err)
		return "", err
	}
	rec.SetWords(1)

	if rec.AcceptWaveform(audioPcmData) != 0 {
		fmt.Println(rec.Result())
	}

	// Unmarshal example for final result
	var jres map[string]interface{}
	json.Unmarshal([]byte(rec.FinalResult()), &jres)
	if text, ok := jres["text"].(string); ok {
		logger.CtxInfo(ctx, "[VoskASR]SpeechToText任务结束,结果：", text, ";耗时：", time.Since(startTime))

		return text, nil
	}
	return "", nil
}
