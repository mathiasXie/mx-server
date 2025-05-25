package audio

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/hraban/opus"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/mathiasXie/gin-web/utils"
)

const FRAME_DURATION_MS = 60

type AudioByte []byte

// 将MP3数据转换为OPUS格式
// mp3Audio: MP3音频数据
// sampleRate: 采样率（如16000）
// channels: 声道数（1为单声道，2为立体声）
// 返回值: 转换后的Opus数据, 音频时长, 错误信息
func AudioToOpusData(ctx context.Context, mp3Audio []byte, sampleRate int, channels int) ([][]byte, float64, error) {

	//tempMapAudioFile, err := os.CreateTemp("", "to_tts_*.mp3")
	tempMapAudioFile, err := os.Create(fmt.Sprintf("./tmp/audio_to_opus_%s_%s.mp3", time.Now().Format("20060102150405"), utils.GetRandomString(10)))

	if err != nil {
		logger.CtxError(ctx, "[AudioUtils]AudioToOpusData无法创建临时音频文件: ", err)
		return nil, 0, err
	}
	err = os.WriteFile(tempMapAudioFile.Name(), mp3Audio, 0644)
	if err != nil {
		logger.CtxError(ctx, "[AudioUtils]AudioToOpusData写入临时音频文件失败:", err)
		return nil, 0, err
	}
	// 从音频文件获取PCM数据 结束时会删除tempMapAudioFile
	pcmData, err := extractPcmFromAudio(tempMapAudioFile.Name())
	if err != nil {
		logger.CtxError(ctx, "[AudioUtils]AudioToOpusData无法从文件提取PCM数据: ", err)
		return nil, 0, err
	}

	// 计算音频时长
	duration := calculateAudioDuration(pcmData, sampleRate, channels)
	// 转换为Opus格式
	opusFrames, err := convertPcmToOpus(pcmData, sampleRate, channels, FRAME_DURATION_MS)
	if err != nil {
		logger.CtxError(ctx, "[AudioUtils]AudioToOpusData转换PCM为Opus失败: ", err)
		return nil, 0, err
	}

	return opusFrames, duration, nil
}

// convertPcmToOpus 将PCM数据转换为Opus格式
// pcmData: PCM音频数据
// sampleRate: 采样率（如16000）
// channels: 声道数（1为单声道，2为立体声）
// frameDurationMs: 帧长度（毫秒，如20、40、60等）
func convertPcmToOpus(pcmData []byte, sampleRate int, channels int, frameDurationMs int) ([][]byte, error) {
	// 创建Opus编码器
	opusEncoder, err := opus.NewEncoder(sampleRate, channels, 2048) // OPUS_APPLICATION_VOIP = 2048
	if err != nil {
		return nil, fmt.Errorf("创建Opus编码器失败: %v", err)
	}

	// 设置编码参数
	if err := opusEncoder.SetBitrate(16000); err != nil { // 16kbps
		return nil, fmt.Errorf("设置比特率失败: %v", err)
	}

	// 计算每帧的采样数
	frameSize := (sampleRate * frameDurationMs) / 1000

	// 将字节数据转换为float32切片
	floatData := make([]float32, len(pcmData)/2)
	for i := 0; i < len(pcmData); i += 2 {
		// 假设PCM数据是16位有符号整数
		sample := int16(binary.LittleEndian.Uint16(pcmData[i:]))
		floatData[i/2] = float32(sample) / 32768.0 // 转换为[-1,1]范围的float32
	}

	// 分帧处理
	var opusFrames [][]byte
	for i := 0; i < len(floatData); i += frameSize {
		end := i + frameSize
		if end > len(floatData) {
			end = len(floatData)
		}

		// 获取当前帧的数据
		frame := floatData[i:end]

		// 如果帧大小不足，用0填充
		if len(frame) < frameSize {
			padding := make([]float32, frameSize-len(frame))
			frame = append(frame, padding...)
		}

		// 编码当前帧
		opusData := make([]byte, 1000) // 预分配足够大的缓冲区
		n, err := opusEncoder.EncodeFloat32(frame, opusData)
		if err != nil {
			return nil, fmt.Errorf("编码帧失败: %v", err)
		}

		opusFrames = append(opusFrames, opusData[:n])
	}

	return opusFrames, nil
}

// convertOpusToPcm 将Opus数据转换为PCM格式
// opusData: Opus编码的音频数据
// sampleRate: 采样率（如16000）
// channels: 声道数（1为单声道，2为立体声）
func ConvertOpusToPcm(opusData []byte, sampleRate int, channels int) ([]byte, error) {
	// 创建Opus解码器
	opusDecoder, err := opus.NewDecoder(sampleRate, channels)
	if err != nil {
		return nil, fmt.Errorf("创建Opus解码器失败: %v", err)
	}

	// 解码Opus数据
	pcmFloat := make([]float32, 960*channels) // 预分配足够大的缓冲区
	n, err := opusDecoder.DecodeFloat32(opusData, pcmFloat)
	if err != nil {
		return nil, fmt.Errorf("解码Opus数据失败: %v", err)
	}

	// 将float32转换为float64
	pcmFloat64 := make([]float64, len(pcmFloat))
	for i, v := range pcmFloat {
		pcmFloat64[i] = float64(v)
	}

	// 应用预加重
	for i := 1; i < len(pcmFloat64); i++ {
		pcmFloat64[i] = pcmFloat64[i] - 0.97*pcmFloat64[i-1]
	}

	// 应用动态范围压缩
	for i := 0; i < len(pcmFloat64); i++ {
		pcmFloat64[i] = math.Min(math.Max(pcmFloat64[i]*0.5, -1.0), 1.0)
	}

	// 将float64转换回float32
	for i, v := range pcmFloat64 {
		pcmFloat[i] = float32(v)
	}

	// 将float32转换为16位整数
	pcmData := make([]byte, n*2)
	for i := 0; i < n; i++ {
		// 将[-1,1]范围的float32转换为16位有符号整数
		intSample := int16(math.Max(-32768, math.Min(32767, float64(pcmFloat[i]*32768))))
		binary.LittleEndian.PutUint16(pcmData[i*2:], uint16(intSample))
	}

	return pcmData, nil
}

func calculateAudioDuration(pcmData []byte, sampleRate int, channels int) float64 {
	// 16位采样
	bytesPerSample := 2

	// 计算总时长（秒）
	return float64(len(pcmData)/1000) / (float64(sampleRate * channels * bytesPerSample))
}

func extractPcmFromAudio(mp3Path string) ([]byte, error) {

	//创建临时PCM文件
	tempPcmFile, err := os.Create(fmt.Sprintf("./tmp/temp_pcm_extract_%s_%s.pcm", time.Now().Format("20060102150405"), utils.GetRandomString(10)))
	if err != nil {
		fmt.Println("无法创建临时PCM文件: ", err)
		return nil, err
	}
	//使用FFmpeg直接将音频转换为PCM
	command := []string{"ffmpeg", "-f", "mp3", "-i", mp3Path, "-f", "s16le", "-acodec", "pcm_s16le", "-y", tempPcmFile.Name()}

	//执行FFmpeg命令
	cmd := exec.Command(command[0], command[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("FFmpeg提取PCM数据失败: ", stderr.String())
		return nil, err
	} else {
		//延迟100ms删除文件
		go func() {
			time.Sleep(500 * time.Millisecond)
			//fmt.Println("将删除临时mp3")
			os.Remove(mp3Path)
			os.Remove(tempPcmFile.Name())
		}()
	}

	// 读取PCM文件内容
	pcmData, err := os.ReadFile(tempPcmFile.Name())
	if err != nil {
		fmt.Println("读取PCM文件内容失败: ", err)
		return nil, err
	}
	return pcmData, nil
}

func CreateAudioDataFromPcm(ctx context.Context, pcmData []byte, format string, sampleRate int, channels int) ([]byte, error) {

	// 将pcm数据写入文件
	tempAudioFile, err := os.Create(fmt.Sprintf("./tmp/to_asr_%s.pcm", time.Now().Format("20060102150405")))
	if err != nil {
		logger.CtxError(ctx, "无法创建临时音频文件: ", err)
		return nil, err
	}
	defer os.Remove(tempAudioFile.Name())
	_, _ = tempAudioFile.Write(pcmData)
	tempAudioFile.Close()

	generrateAudioFile := fmt.Sprintf("./tmp/to_asr_%s.%s", time.Now().Format("20060102150405"), format)
	defer os.Remove(generrateAudioFile)
	//使用ffmpeg将pcm数据转换为指定格式
	command := []string{"ffmpeg", "-f", "s16le", "-ar", strconv.Itoa(sampleRate), "-ac", strconv.Itoa(channels), "-i", tempAudioFile.Name(), generrateAudioFile}

	//执行FFmpeg命令
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	err = cmd.Run()
	if err != nil {
		logger.CtxError(ctx, "FFmpeg提取PCM数据失败: ", err)
		return nil, err
	}

	//读取转换后的音频文件
	audioData, err := os.ReadFile(generrateAudioFile)
	if err != nil {
		logger.CtxError(ctx, "读取转换后的音频文件失败: ", err)
		return nil, err
	}
	return audioData, nil
}
