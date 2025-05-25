package handler

import (
	"encoding/json"
	"fmt"
	"log"
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

		if h.clientListenMode == "auto" {
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
		return nil
	}

	//如果本段有声音，且已经停止了
	if h.clientVoiceStop && h.clientHaveVoice {

		audioData, err := audio_utils.CreateAudioDataFromPcm(h.rpcCtx, h.asrAudio, "mp3", 16000, 1)
		if err != nil {
			h.print(fmt.Sprintf("将asr音频转换为mp3失败:%s", err.Error()), "red")
			logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessage音频转换失败:", err, len(p))
			return err
		}
		text, err := (*h.ASRClient).SpeechToText(h.rpcCtx, &asr_proto.SpeechToTextRequest{
			Provider:  asr_proto.Provider(asr_proto.Provider_value[config.Instance.Provider.ASR]),
			AudioData: audioData,
		})
		if err != nil {
			logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessageASR返回错误:", err)
			return err
		}

		// 清空asr音频
		h.asrAudio = []byte{}
		h.resetVad()

		h.print(fmt.Sprintf("ASR返回结果:%s", text), "green")
		textLen, textResult := utils.RemovePunctuationAndLength(text.Text)
		if textLen > 0 {
			// STT用户说的话
			h.sendTextMessage(textResult, dto.ChatStateStart, dto.ChatTypeSTT)
			h.startToChat(textResult)
		} else {
			textResult = "你说什么？我没听清楚呢"
			h.sendTextMessage(textResult, dto.ChatStateStart, dto.ChatTypeLLM)
			h.sendTTSMessage(dto.ChatStateStart)
			h.sendAudioMessage(textResult)
			h.sendTTSMessage(dto.ChatStateStop)
		}

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

	// 生成TTS音频
	err := h.generateTTSAudio(text, providerName, voiceId, language)
	if err != nil {
		logger.CtxError(h.rpcCtx, "[ChatHandler]sendAudioMessage生成TTS音频失败:", err)
		return err
	}

	return nil
}

// tts音频生成
func (h *ChatHandler) generateTTSAudio(text string, providerName int32, voiceId string, language string) error {

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
			//msg := ttsResp
			if err != nil {
				logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessageTTS返回错误:", err)
				return err
			}
			if len(msg.AudioData) > 0 {

				opusData, _, err := audio_utils.AudioToOpusData(h.rpcCtx, msg.AudioData, 16000, 1)
				if err != nil {
					h.print(fmt.Sprintf("将TTS音频转换为opus失败:%s", err.Error()), "red")
					logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessage音频转换失败:", err)
					return err
				}

				err = h.sendAudio(opusData, true)
				if err != nil {
					continue
				}
			}

			if msg.IsEnd {
				h.print("[ChatHandler]TTS完成一次播报", "cyan")
				break
			}

		}
	}
	return nil
}

// 播放音频
func (h *ChatHandler) sendAudio(audios [][]byte, preBuffer bool) error {

	frameDuration := 60 * time.Millisecond
	startTime := time.Now()
	playPosition := 0 * time.Millisecond
	lastResetTime := time.Now()

	var remainingAudios [][]byte
	if preBuffer {
		preBufferFrames := min(3, len(audios))
		for i := 0; i < preBufferFrames; i++ {
			err := h.conn.WriteMessage(websocket.BinaryMessage, audios[i])
			if err != nil {
				log.Println("预缓冲发送失败:", err)
				return err
			}
		}
		remainingAudios = audios[preBufferFrames:]
	} else {
		remainingAudios = audios
	}

	for _, opusPacket := range remainingAudios {
		if h.clientAbort {
			log.Println("客户端中止，停止发送")
			return nil
		}
		// 每 60 秒重置一次连接超时
		if time.Since(lastResetTime) > time.Minute {
			// if err := h.conn.ResetTimeout(); err != nil {
			// 	log.Println("重置超时失败:", err)
			// 	return nil
			// }
			lastResetTime = time.Now()
		}
		expectedTime := startTime.Add(playPosition)
		delay := time.Until(expectedTime)
		if delay > 0 {
			time.Sleep(delay)
		}
		err := h.conn.WriteMessage(websocket.BinaryMessage, opusPacket)
		if err != nil {
			log.Println("发送失败:", err)
			logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessage发送TTS消息失败:", err)
			return nil
		}

		playPosition += frameDuration

	}
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
		//h.print(fmt.Sprintf("声音检测结果:%t;clientHaveVoice:%t", vadResp.IsActivity, h.clientHaveVoice), "yellow")
		clientHaveVoice = vadResp.IsActivity
		//h.print(fmt.Sprintf("clientHaveVoice:%t", h.clientHaveVoice), "yellow")
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
			//更改clientHaveVoice为true==========
			h.print("检测到声音", "green")
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
