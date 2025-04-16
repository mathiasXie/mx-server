package asr

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

// AliyunASR 实现ASR服务
type AliyunASR struct {
	config Config
	client *http.Client
}

// NewVoskASR
func NewAliyunASR(config Config) (ASRProvider, error) {

	if config.APIKey == "" {
		return nil, fmt.Errorf("microsoft tts requires an API key")
	}
	return &AliyunASR{
		config: config,
		client: &http.Client{},
	}, nil
}

var dialer = websocket.DefaultDialer

const (
	wsURL = "wss://dashscope.aliyuncs.com/api-ws/v1/inference/" // WebSocket服务器地址
)

// SpeechToText 实现ASRProvider接口
func (m *AliyunASR) SpeechToText(ctx context.Context, audioData []byte) (string, error) {

	// 连接WebSocket服务
	conn, err := connectWebSocket(m.config.APIKey)
	if err != nil {
		log.Fatal("连接WebSocket失败：", err)
	}
	defer closeConnection(conn)

	// 发送run-task指令
	taskID, err := sendRunTaskCmd(conn)
	if err != nil {
		log.Fatal("发送run-task指令失败：", err)
	}

	// 发送待识别音频文件流
	if err := sendAudioData(conn, audioData); err != nil {
		log.Fatal("发送音频失败：", err)
	}

	// 发送finish-task指令
	if err := sendFinishTaskCmd(conn, taskID); err != nil {
		log.Fatal("发送finish-task指令失败：", err)
	}

	// 死循环等待响应
	var text string
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("解析服务器消息失败：", err)
			return "", err
		}
		var event Event
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Println("解析事件失败：", err)
			continue
		}
		if event.Header.Event == "result-generated" {
			if event.Payload.Output.Sentence.Text != "" {
				text = event.Payload.Output.Sentence.Text
			}
		} else if event.Header.Event == "task-started" {
			logger.CtxInfo(ctx, "[AliyunASR]SpeechToText任务开始")
		} else if event.Header.Event == "task-finished" || event.Header.Event == "task-failed" {
			logger.CtxInfo(ctx, "[AliyunASR]SpeechToText任务结束,结果：", text)

			break
		}
	}
	return text, nil
}

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

type Output struct {
	Sentence struct {
		BeginTime int64  `json:"begin_time"`
		EndTime   *int64 `json:"end_time"`
		Text      string `json:"text"`
		Words     []struct {
			BeginTime   int64  `json:"begin_time"`
			EndTime     *int64 `json:"end_time"`
			Text        string `json:"text"`
			Punctuation string `json:"punctuation"`
		} `json:"words"`
	} `json:"sentence"`
	Usage interface{} `json:"usage"`
}

type Payload struct {
	TaskGroup  string `json:"task_group"`
	Task       string `json:"task"`
	Function   string `json:"function"`
	Model      string `json:"model"`
	Parameters Params `json:"parameters"`
	// 不使用热词功能时，不要传递resources参数
	// Resources  []Resource `json:"resources"`
	Input  Input  `json:"input"`
	Output Output `json:"output,omitempty"`
}

type Params struct {
	Format                   string   `json:"format"`
	SampleRate               int      `json:"sample_rate"`
	VocabularyID             string   `json:"vocabulary_id"`
	DisfluencyRemovalEnabled bool     `json:"disfluency_removal_enabled"`
	LanguageHints            []string `json:"language_hints"`
}

// 不使用热词功能时，不要传递resources参数
type Resource struct {
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
}

type Input struct {
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
	return conn, err
}

// 发送run-task指令
func sendRunTaskCmd(conn *websocket.Conn) (string, error) {
	runTaskCmd, taskID, err := generateRunTaskCmd()
	if err != nil {
		return "", err
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(runTaskCmd))
	return taskID, err
}

// 生成run-task指令
func generateRunTaskCmd() (string, string, error) {
	taskID := uuid.New().String()
	runTaskCmd := Event{
		Header: Header{
			Action:    "run-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			TaskGroup: "audio",
			Task:      "asr",
			Function:  "recognition",
			Model:     "paraformer-realtime-v2",
			Parameters: Params{
				Format:     "pcm",
				SampleRate: 16000,
			},
			Input: Input{},
		},
	}
	runTaskCmdJSON, err := json.Marshal(runTaskCmd)
	return string(runTaskCmdJSON), taskID, err
}

// 发送音频数据
func sendAudioData(conn *websocket.Conn, audioData []byte) error {

	err := conn.WriteMessage(websocket.BinaryMessage, audioData)
	if err != nil {
		return err
	}
	return nil
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

// 关闭连接
func closeConnection(conn *websocket.Conn) {
	if conn != nil {
		conn.Close()
	}
}
