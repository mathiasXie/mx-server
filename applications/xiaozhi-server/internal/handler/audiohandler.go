package handler

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	asr_proto "github.com/mathiasXie/gin-web/applications/asr-rpc/proto/pb/proto"
	tts_proto "github.com/mathiasXie/gin-web/applications/tts-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/consts"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/mathiasXie/gin-web/utils"
	audio_utils "github.com/mathiasXie/gin-web/utils/audio"
)

func (h *ChatHandler) handlerAudioMessage(p []byte) error {

	if len(p) > 15 {
		pcmData, err := audio_utils.ConvertOpusToPcm(p, 16000, 1)
		if err != nil {
			logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessage音频转换失败:", err, len(p))
			return err
		}
		h.asrAudio = append(h.asrAudio, pcmData...)
		logger.CtxInfo(h.rpcCtx, "[ChatHandler]handlerAudioMessage音频转换成功:", len(pcmData))
	} else {
		logger.CtxInfo(h.rpcCtx, "[ChatHandler]handlerAudioMessage音频长度小于15:", len(p))
	}
	//如果本段有声音，且已经停止了
	if h.clientVoiceStop {
		audioData, err := audio_utils.CreateAudioDataFromPcm(h.rpcCtx, h.asrAudio, "wav", 16000, 1)
		if err != nil {
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

	}
	return nil
}

func (h *ChatHandler) sendAudioMessage(text string) error {

	//获取TTS配置
	providerName, ok := tts_proto.Provider_value[h.userInfo.Role.TTS]
	if !ok {
		logger.CtxError(h.rpcCtx, "[ChatHandler]sendAudioMessage用户的TTS配置错误:", h.userInfo.Role.TTS)
		return consts.RetTTSConfigError
	}
	provider := tts_proto.Provider(providerName)
	h.sendTextMessage(text, dto.ChatStateSentenceStart, dto.ChatTypeTTS)

	// 发送一段语音给客户端
	ttsResp, err := (*h.TTSClient).TextToSpeechStream(h.rpcCtx, &tts_proto.TextToSpeechRequest{
		Provider: provider,
		VoiceId:  h.userInfo.Role.TTSVoiceId,
		Language: h.userInfo.Role.Language,
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
					logger.CtxError(h.rpcCtx, "[ChatHandler]handlerAudioMessage音频转换失败:", err)
					return err
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
