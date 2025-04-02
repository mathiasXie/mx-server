package server

import (
	"context"
	"fmt"

	tts "github.com/mathiasXie/gin-web/applications/tts-rpc/internal"
	"github.com/mathiasXie/gin-web/applications/tts-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/consts"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	proto.UnimplementedTTSServiceServer
	ttsConfigs map[proto.Provider]tts.Config
}

func NewServer() *Server {

	ttsConfigs := map[proto.Provider]tts.Config{
		proto.Provider_MICROSOFT: {
			APIKey:   config.Instance.TTS.Microsoft.APIKey,
			Endpoint: config.Instance.TTS.Microsoft.Endpoint,
		},
		proto.Provider_DOUBAO: {
			APIKey:   "your_api_key",
			Endpoint: "https://api.tts.com/v1",
		},
	}

	return &Server{
		ttsConfigs: ttsConfigs,
	}
}

func (s *Server) getTTSProvider(ctx context.Context, provider proto.Provider) (tts.TTSProvider, context.Context, error) {
	ttsConfig, ok := s.ttsConfigs[provider]
	if !ok {
		logger.CtxError(ctx, "failed to create tts provider: %v", provider)
		return nil, ctx, fmt.Errorf("failed to create tts provider: %v", provider)
	}
	var ttsProvider tts.TTSProvider
	var err error
	switch provider {
	case proto.Provider_MICROSOFT:
		ttsProvider, err = tts.NewMicrosoftTTS(ttsConfig)
	case proto.Provider_DOUBAO:
		ttsProvider, err = tts.NewDoubaoTTS(ttsConfig)
	}
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to create tts provider: %v", err)
	}
	// 从上下文中获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// 打印元数据
		trace_id, ok := md["trace_id"]
		if ok {
			trace_id_str := trace_id[0]
			ctx = context.WithValue(ctx, consts.LogID, trace_id_str)
		}
	}
	return ttsProvider, ctx, nil
}

func (s *Server) TextToSpeech(ctx context.Context, req *proto.TextToSpeechRequest) (*proto.TextToSpeechResponse, error) {

	ttsProvider, ctx, err := s.getTTSProvider(ctx, req.Provider)
	if err != nil {
		logger.CtxError(ctx, "failed to get tts provider: %v", err)
		return nil, fmt.Errorf("failed to get tts provider: %v", err)
	}

	audioData, format, sampleRate, channels, err := ttsProvider.TextToSpeech(
		ctx,
		req.Text,
		req.Language,
		req.VoiceId,
		float32(req.Speed),
		float32(req.Pitch),
	)
	if err != nil {
		logger.CtxError(ctx, "failed to convert text to speech: %v", err)
		return nil, fmt.Errorf("failed to convert text to speech: %v", err)
	}
	logger.CtxInfo(ctx, "req: %v", req)
	return &proto.TextToSpeechResponse{
		AudioData:  audioData,
		Format:     format,
		SampleRate: sampleRate,
		Channels:   channels,
	}, nil
}

func (s *Server) VoicesList(ctx context.Context, req *proto.VoicesListRequest) (*proto.VoicesListResponse, error) {

	ttsProvider, ctx, err := s.getTTSProvider(ctx, req.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get tts provider: %v", err)
	}

	voices, err := ttsProvider.VoicesList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get voices list: %v", err)
	}

	protoVoices := make([]*proto.Voice, len(voices))
	for i, voice := range voices {
		protoVoices[i] = &proto.Voice{
			VoiceId:   voice.VoiceID,
			VoiceName: voice.VoiceName,
		}
	}
	return &proto.VoicesListResponse{Voices: protoVoices}, nil
}
