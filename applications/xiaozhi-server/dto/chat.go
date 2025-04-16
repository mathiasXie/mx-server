package dto

// {type: 'tts', state: 'sentence_start', session_id: '09fc016d-3425-462d-a92c-782ef53dacca', text: '1. 将鸡蛋打入碗中，加入适量盐，搅拌均匀'}

// type: tts/llm/stt/hello
// state: start/sentence_start/sentence_end/stop
// type用枚举类型
type ChatType string
type ChatState string

const (
	ChatTypeTTS    ChatType = "tts"
	ChatTypeLLM    ChatType = "llm"
	ChatTypeSTT    ChatType = "stt"
	ChatTypeHello  ChatType = "hello"
	ChatTypeListen ChatType = "listen"
)

const (
	ChatStateStart         ChatState = "start"
	ChatStateSentenceStart ChatState = "sentence_start"
	ChatStateSentenceEnd   ChatState = "sentence_end"
	ChatStateStop          ChatState = "stop"
	ChatStateAuto          ChatState = "auto"
	ChatStateDetect        ChatState = "detect"
)

type ChatResponse struct {
	Type      ChatType  `json:"type"`
	State     ChatState `json:"state"`
	SessionID string    `json:"session_id"`
	Text      string    `json:"text"`
	Emotion   string    `json:"emotion"`
}

// {"type":"hello","device_id":"web_test_device","device_name":"Web测试设备","device_mac":"00:11:22:33:44:55","token":"your-token1"}
type ChatRequest struct {
	Type       ChatType  `json:"type"`
	DeviceID   string    `json:"device_id"`
	DeviceName string    `json:"device_name"`
	DeviceMac  string    `json:"device_mac"`
	Token      string    `json:"token"`
	Mode       string    `json:"mode"`
	State      ChatState `json:"state"`
}
