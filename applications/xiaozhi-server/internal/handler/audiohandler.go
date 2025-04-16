package handler

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

/*
async def handleAudioMessage(conn, audio):
    if not conn.asr_server_receive:
        logger.bind(tag=TAG).debug(f"前期数据处理中，暂停接收")
        return
    if conn.client_listen_mode == "auto":
        have_voice = conn.vad.is_vad(conn, audio)
    else:
        have_voice = conn.client_have_voice

    # 如果本次没有声音，本段也没声音，就把声音丢弃了
    if have_voice == False and conn.client_have_voice == False:
        await no_voice_close_connect(conn)
        conn.asr_audio.append(audio)
        conn.asr_audio = conn.asr_audio[
            -10:
        ]  # 保留最新的10帧音频内容，解决ASR句首丢字问题
        return
    conn.client_no_voice_last_time = 0.0
    conn.asr_audio.append(audio)
    # 如果本段有声音，且已经停止了
    if conn.client_voice_stop:
        conn.client_abort = False
        conn.asr_server_receive = False
        # 音频太短了，无法识别
        if len(conn.asr_audio) < 15:
            conn.asr_server_receive = True
        else:
            text, file_path = await conn.asr.speech_to_text(
                conn.asr_audio, conn.session_id
            )
            logger.bind(tag=TAG).info(f"识别文本: {text}")
            text_len, _ = remove_punctuation_and_length(text)
            if text_len > 0:
                await startToChat(conn, text)
            else:
                conn.asr_server_receive = True
        conn.asr_audio.clear()
        conn.reset_vad_states()
*/

func (h *ChatHandler) handlerAudioMessage(rpcCtx context.Context, p []byte, conn *websocket.Conn) error {

	//如果本次没有声音，本段也没声音，就把声音丢弃了

	h.asrAudio = append(h.asrAudio, p...)

	//如果本段有声音，且已经停止了
	if h.clientVoiceStop {

		//将音频写入文件
		tempAudioFile, err := os.Create("./tmp/to_tts_.opus")
		if err != nil {
			fmt.Println("无法创建临时音频文件: ", err)
			return err
		}
		err = os.WriteFile(tempAudioFile.Name(), h.asrAudio, 0644)
		if err != nil {
			log.Println("Error writing message:", err)
			return err
		}
		fmt.Println(tempAudioFile.Name())
	}
	return nil
}
