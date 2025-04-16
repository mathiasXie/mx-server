package audio

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"os/exec"

	"github.com/hraban/opus"
)

const FRAME_DURATION_MS = 60

type AudioByte []byte

// ConvertMp3ToOpusBytes 将MP3数据转换为OPUS格式
func AudioToOpusData(mp3Path string, sampleRate int, channels int) ([]AudioByte, float64, error) {
	// 从音频文件获取PCM数据
	pcmData, err := extractPcmFromAudio(mp3Path)
	if err != nil {
		fmt.Println("无法从文件提取PCM数据: ", mp3Path)
		return nil, 0, err
	}

	// 计算音频时长
	duration := calculateAudioDuration(pcmData, sampleRate, channels)
	// 转换为Opus格式
	opusFrames, err := convertPcmToOpus(pcmData, sampleRate, channels, FRAME_DURATION_MS)
	if err != nil {
		fmt.Println("转换PCM为Opus失败: ", err)
		return nil, 0, err
	}

	return opusFrames, duration, nil
}

// convertPcmToOpus 将PCM数据转换为Opus格式
// pcmData: PCM音频数据
// sampleRate: 采样率（如16000）
// channels: 声道数（1为单声道，2为立体声）
// frameDurationMs: 帧长度（毫秒，如20、40、60等）
func convertPcmToOpus(pcmData []byte, sampleRate int, channels int, frameDurationMs int) ([]AudioByte, error) {
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
	var opusFrames []AudioByte
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
func convertOpusToPcm(opusData []byte, sampleRate int, channels int) ([]byte, error) {
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
	tempPcmFile, err := os.CreateTemp("", "temp_pcm_extract_*.pcm")
	if err != nil {
		fmt.Println("无法创建临时PCM文件: ", err)
		return nil, err
	}
	defer os.Remove(tempPcmFile.Name())

	//使用FFmpeg直接将音频转换为PCM
	command := []string{"ffmpeg", "-i", mp3Path, "-f", "s16le", "-acodec", "pcm_s16le", "-y", tempPcmFile.Name()}

	//执行FFmpeg命令
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("FFmpeg提取PCM数据失败: ", err)
		return nil, err
	}

	// 读取PCM文件内容
	pcmData, err := os.ReadFile(tempPcmFile.Name())
	if err != nil {
		fmt.Println("读取PCM文件内容失败: ", err)
		return nil, err
	}
	return pcmData, nil
}
