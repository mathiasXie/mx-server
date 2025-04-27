package asr

import (
	"context"
	"fmt"
	"math"

	"github.com/mathiasXie/gin-web/applications/asr-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/streamer45/silero-vad-go/speech"
)

// 从文件中读取音频样本
func readSamplesFromPCMByte(data []byte) []float32 {

	samples := make([]float32, 0, len(data)/4)
	for i := 0; i < len(data); i += 4 {
		samples = append(samples, math.Float32frombits(uint32(data[i])|uint32(data[i+1])<<8|uint32(data[i+2])<<16|uint32(data[i+3])<<24))
	}
	return samples
}

func readInt16PCMAsFloat32(data []byte) []float32 {
	samples := make([]float32, 0, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		// 16-bit signed little-endian
		sample := int16(data[i]) | int16(data[i+1])<<8
		// 归一化到 -1.0 ~ +1.0
		samples = append(samples, float32(sample)/32768.0)
	}
	return samples
}

func SileroVadDetect(ctx context.Context, audioData []byte) (*proto.DetectVADResponse, error) {
	isActivity := false
	// 配置检测器
	cfg := speech.DetectorConfig{
		ModelPath:            fmt.Sprintf("./models/%s", "silero_vad.onnx"),
		SampleRate:           16000,
		Threshold:            0.8, //数值越大越保守
		MinSilenceDurationMs: 300,
		SpeechPadMs:          50,
	}
	// 创建检测器实例
	sd, err := speech.NewDetector(cfg)
	if err != nil {
		logger.CtxError(ctx, "Failed to create detector", err)
		return nil, err
	}
	defer func() {
		if err := sd.Destroy(); err != nil {
			logger.CtxError(ctx, "Failed to detect speech", err)
		}
	}()

	// 读取音频样本
	samples := readInt16PCMAsFloat32(audioData)
	if samples != nil {
		// 进行语音检测
		segments, err := sd.Detect(samples)
		if err != nil {
			//log.Fatal("Failed to detect speech", "error", err)
			logger.CtxError(ctx, "Failed to detect speech", err)
			return nil, err
		}
		isActivity = len(segments) > 0
		// fmt.Println("segments", len(segments))

		// // 输出检测结果
		// for _, segment := range segments {
		// 	fmt.Printf("Speech start at: %.3f s, end at: %.3f s\n", segment.SpeechStartAt, segment.SpeechEndAt)
		// }
	}
	return &proto.DetectVADResponse{
		IsActivity: isActivity,
	}, nil
}
