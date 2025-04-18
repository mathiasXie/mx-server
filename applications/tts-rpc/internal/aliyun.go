package tts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

const (
	wsURL  = "wss://dashscope.aliyuncs.com/api-ws/v1/inference/" // WebSocket服务器地址
	format = "mp3"
)

// AliyunTTS 实现微软的TTS服务
type AliyunTTS struct {
	config Config
	client *http.Client
}

// NewAliyunTTS
func NewAliyunTTS(config Config) (TTSProvider, error) {

	if config.APIKey == "" {
		return nil, fmt.Errorf("aliyun tts requires an API key")
	}
	return &AliyunTTS{
		config: config,
		client: &http.Client{},
	}, nil
}

// TextToSpeech 实现TTSProvider接口
func (m *AliyunTTS) TextToSpeech(ctx context.Context, text string, language string, voiceID string, speed float32, pitch float32) ([]byte, string, int32, int32, error) {

	// 连接WebSocket服务
	conn, err := connectWebSocket(m.config.APIKey)
	if err != nil {
		fmt.Println("连接WebSocket失败：", err)
		logger.CtxError(ctx, "[AliyunTTS][TextToSpeech]连接WebSocket失败：%v", err)
		return nil, "", 0, 0, err
	}
	defer closeConnection(conn)

	// 发送run-task指令
	taskID, err := sendRunTaskCmd(conn, voiceID)
	if err != nil {
		fmt.Println("发送run-task指令失败：", err)
		logger.CtxError(ctx, "[AliyunTTS][TextToSpeech]发送run-task指令失败：%v", err)
		return nil, "", 0, 0, err
	}

	// 发送待合成文本
	if err := sendContinueTaskCmd(conn, taskID, text); err != nil {
		fmt.Println("发送待合成文本失败：", err)
		logger.CtxError(ctx, "[AliyunTTS][TextToSpeech]发送待合成文本失败：%v", err)
		return nil, "", 0, 0, err
	}

	// 发送finish-task指令
	if err := sendFinishTaskCmd(conn, taskID); err != nil {
		fmt.Println("发送finish-task指令失败：", err)
		logger.CtxError(ctx, "发送finish-task指令失败：%v", err)
		return nil, "", 0, 0, err
	}

	// 循环等待响应
	var audio []byte
	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("解析服务器消息失败：", err)
			logger.CtxError(ctx, "[AliyunTTS][TextToSpeech]解析服务器消息失败：%v", err)
			return nil, "", 0, 0, err
		}

		if msgType == websocket.BinaryMessage {
			audio = append(audio, message...)
		} else {
			// 处理文本消息
			var event Event
			err = json.Unmarshal(message, &event)
			if err != nil {
				fmt.Println("解析事件失败：", err)
				logger.CtxError(ctx, "[AliyunTTS][TextToSpeech]解析事件失败：%v", err)
				continue
			}

			if event.Header.Event == "task-finished" {
				logger.CtxInfo(ctx, "[AliyunTTS][TextToSpeech]任务完成")
				break
			} else if event.Header.Event == "task-started" {
				logger.CtxInfo(ctx, "[AliyunTTS]TextToSpeech任务开始,报文：", text)
			}
		}
	}
	return audio, format, 22050, 1, nil
}

func (m *AliyunTTS) TextToSpeechStream(ctx context.Context, text string, language string, voiceID string, respChan chan<- TTSStreamResponse) error {
	start := time.Now()
	// 连接WebSocket服务
	conn, err := connectWebSocket(m.config.APIKey)
	if err != nil {
		fmt.Println("连接WebSocket失败：", err)
		logger.CtxError(ctx, "[AliyunTTS][TextToSpeech]连接WebSocket失败：%v", err)
		return err
	}
	defer closeConnection(conn)

	// 发送run-task指令
	taskID, err := sendRunTaskCmd(conn, voiceID)
	if err != nil {
		fmt.Println("发送run-task指令失败：", err)
		logger.CtxError(ctx, "[AliyunTTS][TextToSpeech]发送run-task指令失败：%v", err)
		return err
	}

	// 发送待合成文本
	if err := sendContinueTaskCmd(conn, taskID, text); err != nil {
		fmt.Println("发送待合成文本失败：", err)
		logger.CtxError(ctx, "[AliyunTTS][TextToSpeech]发送待合成文本失败：%v", err)
		return err
	}

	// 发送finish-task指令
	if err := sendFinishTaskCmd(conn, taskID); err != nil {
		fmt.Println("发送finish-task指令失败：", err)
		logger.CtxError(ctx, "发送finish-task指令失败：%v", err)
		return err
	}

	// 循环等待响应
	var audio []byte
	buferLenght := 10240
	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("解析服务器消息失败：", err)
			logger.CtxError(ctx, "[AliyunTTS][TextToSpeech]解析服务器消息失败：%v", err)
			return err
		}

		if msgType == websocket.BinaryMessage {
			audio = append(audio, message...)
			for len(audio) >= buferLenght {
				sendAudio := audio[:buferLenght]
				respChan <- TTSStreamResponse{
					Audio:  sendAudio,
					Sample: 16000,
					Format: format,
				}
				audio = audio[buferLenght:]
			}
		} else {
			// 处理文本消息
			var event Event
			err = json.Unmarshal(message, &event)
			if err != nil {
				fmt.Println("解析事件失败：", err)
				logger.CtxError(ctx, "[AliyunTTS][TextToSpeech]解析事件失败：%v", err)
				continue
			}
			if event.Header.Event == "task-finished" {
				// 处理剩余不足 10240 字节的数据
				if len(audio) > 0 {
					respChan <- TTSStreamResponse{
						Audio:  audio,
						Sample: 16000,
						Format: format,
					}
				}
				respChan <- TTSStreamResponse{
					IsEnd: true,
				}
				logger.CtxInfo(ctx, "[AliyunTTS][TextToSpeech]TextToSpeechStream输出结束:", text, ",耗时：", time.Since(start).Milliseconds())
				break
			}
		}
	}
	return nil
}

func (m *AliyunTTS) VoicesList(ctx context.Context) ([]Voices, error) {
	voices := make([]Voices, 0)

	return voices, nil
}

var dialer = websocket.DefaultDialer

// 定义结构体来表示JSON数据
type Header struct {
	Action       string                 `json:"action"`
	TaskID       string                 `json:"task_id"`
	Streaming    string                 `json:"streaming"`
	Event        string                 `json:"event"`
	ErrorCode    string                 `json:"error_code,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	Attributes   map[string]interface{} `json:"attributes"`
}

type Payload struct {
	TaskGroup  string     `json:"task_group"`
	Task       string     `json:"task"`
	Function   string     `json:"function"`
	Model      string     `json:"model"`
	Parameters Params     `json:"parameters"`
	Resources  []Resource `json:"resources"`
	Input      Input      `json:"input"`
}

type Params struct {
	TextType   string `json:"text_type"`
	Voice      string `json:"voice"`
	Format     string `json:"format"`
	SampleRate int    `json:"sample_rate"`
	Volume     int    `json:"volume"`
	Rate       int    `json:"rate"`
	Pitch      int    `json:"pitch"`
}

type Resource struct {
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
}

type Input struct {
	Text string `json:"text"`
}

type Event struct {
	Header  Header  `json:"header"`
	Payload Payload `json:"payload"`
}

// 连接WebSocket服务
func connectWebSocket(apiKey string) (*websocket.Conn, error) {
	header := make(http.Header)
	header.Add("X-DashScope-DataInspection", "enable")
	header.Add("Authorization", fmt.Sprintf("bearer %s", apiKey))
	conn, _, err := dialer.Dial(wsURL, header)
	if err != nil {
		fmt.Println("连接WebSocket失败：", err)
		return nil, err
	}
	return conn, nil
}

// 发送run-task指令
func sendRunTaskCmd(conn *websocket.Conn, voiceID string) (string, error) {
	runTaskCmd, taskID, err := generateRunTaskCmd(voiceID)
	if err != nil {
		return "", err
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(runTaskCmd))
	return taskID, err
}

// 生成run-task指令
func generateRunTaskCmd(voiceID string) (string, string, error) {
	taskID := uuid.New().String()
	runTaskCmd := Event{
		Header: Header{
			Action:    "run-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			TaskGroup: "audio",
			Task:      "tts",
			Function:  "SpeechSynthesizer",
			Model:     "cosyvoice-v1",
			Parameters: Params{
				TextType:   "PlainText",
				Voice:      voiceID,
				Format:     format,
				SampleRate: 16000,
				Volume:     50,
				Rate:       1,
				Pitch:      1,
			},
			Input: Input{},
		},
	}
	runTaskCmdJSON, err := json.Marshal(runTaskCmd)
	return string(runTaskCmdJSON), taskID, err
}

// 发送待合成文本
func sendContinueTaskCmd(conn *websocket.Conn, taskID string, text string) error {
	texts := []string{text}

	for _, text := range texts {
		runTaskCmd, err := generateContinueTaskCmd(text, taskID)
		if err != nil {
			return err
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte(runTaskCmd))
		if err != nil {
			return err
		}
	}

	return nil
}

// 生成continue-task指令
func generateContinueTaskCmd(text string, taskID string) (string, error) {
	runTaskCmd := Event{
		Header: Header{
			Action:    "continue-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			Input: Input{
				Text: text,
			},
		},
	}
	runTaskCmdJSON, err := json.Marshal(runTaskCmd)
	return string(runTaskCmdJSON), err
}

// 关闭连接
func closeConnection(conn *websocket.Conn) {
	if conn != nil {
		conn.Close()
	}
}

// 发送finish-task指令
func sendFinishTaskCmd(conn *websocket.Conn, taskID string) error {
	finishTaskCmd, err := generateFinishTaskCmd(taskID)
	if err != nil {
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(finishTaskCmd))
	return err
}

// 生成finish-task指令
func generateFinishTaskCmd(taskID string) (string, error) {
	finishTaskCmd := Event{
		Header: Header{
			Action:    "finish-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			Input: Input{},
		},
	}
	finishTaskCmdJSON, err := json.Marshal(finishTaskCmd)
	return string(finishTaskCmdJSON), err
}
