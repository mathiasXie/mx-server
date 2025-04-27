package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	asr_proto "github.com/mathiasXie/gin-web/applications/asr-rpc/proto/pb/proto"
	tts_proto "github.com/mathiasXie/gin-web/applications/tts-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/consts"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/mathiasXie/gin-web/utils"
	audio_utils "github.com/mathiasXie/gin-web/utils/audio"
)

func (h *ChatHandler) handlerAudioMessage(p []byte) error {

	var haveVoice bool
	if len(p) > 15 {
		pcmData, err := audio_utils.ConvertOpusToPcm(p, 16000, 1)
		if err != nil {
			h.print(fmt.Sprintf("将客户端音频转换为pcm失败:%s", err.Error()), "red")
			logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessage音频转换失败:", err, len(p))
			return err
		}

		if h.clientListenMode == "auto" || 1 == 1 {
			// 处理vad检测
			haveVoice, err = h.handleVad(pcmData)
			if err != nil {
				logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessage vad检测失败:", err)
				return err
			}
		} else {
			haveVoice = h.clientHaveVoice
		}
		//如果本次没有声音，本段也没声音，就把声音丢弃了
		if !haveVoice && !h.clientHaveVoice {
			h.asrAudio = append(h.asrAudio, pcmData...)
			//保留最新的10帧音频内容，解决ASR句首丢字问题
			h.asrAudio = h.asrAudio[len(h.asrAudio)-10:]
			return nil
		}
		h.asrAudio = append(h.asrAudio, pcmData...)

	} else {
		logger.CtxInfo(h.rpcCtx, "[ChatHandler]handlerAudioMessage音频长度小于15:", len(p))
	}

	//如果本段有声音，且已经停止了
	if h.clientVoiceStop && h.clientHaveVoice {

		audioData, err := audio_utils.CreateAudioDataFromPcm(h.rpcCtx, h.asrAudio, "wav", 16000, 1)
		if err != nil {
			h.print(fmt.Sprintf("将asr音频转换为wav失败:%s", err.Error()), "red")
			logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessage音频转换失败:", err, len(p))
			return err
		}
		text, err := (*h.ASRClient).SpeechToText(h.rpcCtx, &asr_proto.SpeechToTextRequest{
			Provider:  asr_proto.Provider_VOSK,
			AudioData: audioData,
		})
		if err != nil {
			logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessageASR返回错误:", err)
			return err
		}
		textLen, textResult := utils.RemovePunctuationAndLength(text.Text)
		if textLen > 0 {
			// 发送TTS消息
			h.sendTextMessage(textResult, dto.ChatStateStart, dto.ChatTypeSTT)
			h.startToChat(textResult)

		} else {
			textResult = "你说什么？我没听清楚呢"
			h.sendTextMessage(textResult, dto.ChatStateStart, dto.ChatTypeLLM)
			h.sendAudioMessage(textResult)
		}
		// 清空asr音频
		h.asrAudio = []byte{}
		h.resetVad()

	}
	return nil
}

func (h *ChatHandler) sendAudioMessage(text string) error {

	var providerName int32
	var voiceId string
	var language string
	// 用户没有配置角色，可能是没有绑定设备，使用系统默认声音id
	if h.userInfo.Role == nil {
		providerName = tts_proto.Provider_value[config.Instance.Provider.TTS]
		voiceId = config.Instance.Provider.DefaultVoice
	} else {
		//获取TTS配置
		var ok bool
		providerName, ok = tts_proto.Provider_value[h.userInfo.Role.TTS]
		if !ok {
			logger.CtxError(h.rpcCtx, "[ChatHandler]sendAudioMessage用户的TTS配置错误:", h.userInfo.Role.TTS)
			return consts.RetTTSConfigError
		}
		voiceId = h.userInfo.Role.TTSVoiceId
		language = h.userInfo.Role.Language
	}
	// 服务器发送语音段开始消息
	h.sendTextMessage(text, dto.ChatStateSentenceStart, dto.ChatTypeTTS)

	// 发送一段语音给客户端
	ttsResp, err := (*h.TTSClient).TextToSpeechStream(h.rpcCtx, &tts_proto.TextToSpeechRequest{
		Provider: tts_proto.Provider(providerName),
		VoiceId:  voiceId,
		Language: language,
		Text:     text,
	})
	if err != nil {
		logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessage连接TTS失败:", err)
		return err
	} else {
		for {
			msg, err := ttsResp.Recv()
			if err != nil {
				logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessageTTS返回错误:", err)
				return err
			}
			if len(msg.AudioData) > 0 {
				opusData, _, err := audio_utils.AudioToOpusData(h.rpcCtx, msg.AudioData, 16000, 1)
				if err != nil {
					h.print(fmt.Sprintf("将TTS音频转换为opus失败:%s", err.Error()), "red")
					logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessage音频转换失败:", err)
					continue
				}
				for _, data := range opusData {
					err := h.conn.WriteMessage(websocket.BinaryMessage, data)
					if err != nil {
						logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessage发送TTS消息失败:", err)
						break
					}
				}
			}

			if msg.IsEnd {
				h.print("[ChatHandler]TTS完成一次播报", "cyan")
				break
			}
		}
	}
	// 服务器发送语音段结束消息
	h.sendTextMessage(text, dto.ChatStateSentenceEnd, dto.ChatTypeTTS)

	return nil
}

func (h *ChatHandler) sendTTSMessage(state dto.ChatState) error {
	chatResponse := dto.ChatResponse{
		Type:      dto.ChatTypeTTS,
		State:     state,
		SessionID: h.sessionID,
	}
	resp, _ := json.Marshal(chatResponse)
	err := h.conn.WriteMessage(websocket.TextMessage, resp)
	if err != nil {
		logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessage发送TTS消息失败:", err)
		return err
	}
	return nil
}

func (h *ChatHandler) handleVad(pcm []byte) (bool, error) {

	// 加入缓冲区
	h.vadAudioBuffer = append(h.vadAudioBuffer, pcm...)

	var clientHaveVoice bool

	// 每次取 512 个采样点（1024字节） 由于golang的vad包不支持流式检测，语音太短，检测不出来，所以乘10
	frameSize := 512 * 2 * 10
	for len(h.vadAudioBuffer) >= frameSize {
		chunk := h.vadAudioBuffer[:frameSize]
		h.vadAudioBuffer = h.vadAudioBuffer[frameSize:]

		// 检测是否有语音
		vadResp, err := (*h.ASRClient).DetectVAD(h.rpcCtx, &asr_proto.DetectVADRequest{
			AudioData: chunk,
		})
		if err != nil {
			h.print(err.Error(), "red")
			return false, err
		}
		//h.print(fmt.Sprintf("声音检测结果:%t", vadResp.IsActivity), "yellow")
		clientHaveVoice = vadResp.IsActivity

		// 如果之前有声音但现在没声音，计算静默时间
		if h.clientHaveVoice && !clientHaveVoice {
			stopMs := time.Now().UnixMilli() - h.vadLastHaveVoiceTime
			if stopMs >= h.vadSilenceThreshold {
				h.print("vad检测到一段语音结束", "green")
				h.clientVoiceStop = true
			}
		}

		if clientHaveVoice {
			h.clientHaveVoice = true
			h.vadLastHaveVoiceTime = time.Now().UnixMilli()
		}
	}

	return clientHaveVoice, nil
}

func (h *ChatHandler) resetVad() {
	h.vadLastHaveVoiceTime = 0
	h.vadAudioBuffer = []byte{}
	h.clientHaveVoice = false
	h.clientVoiceStop = false
}
