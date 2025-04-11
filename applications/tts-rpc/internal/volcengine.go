package tts

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// VolcEngineTTS 实现豆包的TTS服务
type VolcEngineTTS struct {
	config Config
	client *http.Client
}

var defaultHeader = []byte{0x11, 0x10, 0x11, 0x00}

type synResp struct {
	Audio  []byte
	IsLast bool
}

const (
	optQuery  string = "query"
	optSubmit string = "submit"
)

// NewVolcEngineTTS 创建新的火山引擎TTS实例
func NewVolcEngineTTS(config Config) (TTSProvider, error) {
	if config.APIID == "" {
		return nil, fmt.Errorf("VolcEngine tts requires an API id")
	}
	return &VolcEngineTTS{
		config: config,
		client: &http.Client{},
	}, nil
}

// TextToSpeech 实现TTSProvider接口
func (d *VolcEngineTTS) TextToSpeech(ctx context.Context, text string, language string, voiceID string, speed float32, pitch float32) ([]byte, string, int32, int32, error) {
	input := setupInput(text, voiceID, optSubmit, d.config.APIID)
	input = gzipCompress(input)
	payloadSize := len(input)
	payloadArr := make([]byte, 4)
	binary.BigEndian.PutUint32(payloadArr, uint32(payloadSize))
	clientRequest := make([]byte, len(defaultHeader))
	copy(clientRequest, defaultHeader)
	clientRequest = append(clientRequest, payloadArr...)
	clientRequest = append(clientRequest, input...)

	var header = http.Header{"Authorization": []string{fmt.Sprintf("Bearer;%s", d.config.Token)}}

	c, _, err := websocket.DefaultDialer.Dial(d.config.Endpoint, header)
	if err != nil {
		fmt.Println("dial err:", err)
		return nil, "", 0, 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer c.Close()
	err = c.WriteMessage(websocket.BinaryMessage, clientRequest)
	if err != nil {
		fmt.Println("write message fail, err:", err.Error())
		return nil, "", 0, 0, fmt.Errorf("failed to send request: %v", err)
	}
	var audio []byte
	for {
		var message []byte
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read message fail, err:", err.Error())
			break
		}
		resp, err := parseResponse(message)
		if err != nil {
			fmt.Println("parse response fail, err:", err.Error())
			break
		}
		audio = append(audio, resp.Audio...)
		if resp.IsLast {
			break
		}
	}
	if err != nil {
		fmt.Println("stream synthesis fail, err:", err.Error())
		return nil, "", 0, 0, fmt.Errorf("failed to send request: %v", err)
	}
	return audio, "mp3", 16000, 1, nil
}

func (m *VolcEngineTTS) VoicesList(ctx context.Context) ([]Voices, error) {

	voices := make([]Voices, 0)

	return voices, nil
}
func setupInput(text, voiceType, opt, appid string) []byte {
	reqID := uuid.Must(uuid.NewV4(), nil).String()
	params := make(map[string]map[string]interface{})
	params["app"] = make(map[string]interface{})
	//平台上查看具体appid
	params["app"]["appid"] = appid
	params["app"]["token"] = "access_token"
	//平台上查看具体集群名称
	params["app"]["cluster"] = "volcano_tts"
	params["user"] = make(map[string]interface{})
	params["user"]["uid"] = "uid"
	params["audio"] = make(map[string]interface{})
	params["audio"]["voice_type"] = voiceType
	params["audio"]["encoding"] = "opus"
	params["audio"]["speed_ratio"] = 1.0
	params["audio"]["volume_ratio"] = 1.0
	params["audio"]["pitch_ratio"] = 1.0
	params["request"] = make(map[string]interface{})
	params["request"]["reqid"] = reqID
	params["request"]["text"] = text
	params["request"]["text_type"] = "plain"
	params["request"]["operation"] = opt
	resStr, _ := json.Marshal(params)
	return resStr
}

func gzipCompress(input []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(input)
	w.Close()
	return b.Bytes()
}

func gzipDecompress(input []byte) []byte {
	b := bytes.NewBuffer(input)
	r, _ := gzip.NewReader(b)
	out, _ := ioutil.ReadAll(r)
	r.Close()
	return out
}

func parseResponse(res []byte) (resp synResp, err error) {
	headSize := res[0] & 0x0f
	messageType := res[1] >> 4
	messageTypeSpecificFlags := res[1] & 0x0f
	messageCompression := res[2] & 0x0f
	payload := res[headSize*4:]

	// audio-only server response
	if messageType == 0xb {
		// no sequence number as ACK
		if messageTypeSpecificFlags == 0 {
		} else {
			sequenceNumber := int32(binary.BigEndian.Uint32(payload[0:4]))
			payload = payload[8:]
			resp.Audio = append(resp.Audio, payload...)

			if sequenceNumber < 0 {
				resp.IsLast = true
			}
		}
	} else if messageType == 0xf {
		errMsg := payload[8:]
		if messageCompression == 1 {
			errMsg = gzipDecompress(errMsg)
		}

		err = errors.New(string(errMsg))
		return
	} else if messageType == 0xc {
		//msgSize := int32(binary.BigEndian.Uint32(payload[0:4]))
		payload = payload[4:]
		if messageCompression == 1 {
			payload = gzipDecompress(payload)
		}
	} else {
		err = errors.New("wrong message type")
		return
	}
	return
}
