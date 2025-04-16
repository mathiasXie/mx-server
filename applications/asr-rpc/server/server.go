package server

import (
	"context"
	"fmt"

	asr "github.com/mathiasXie/gin-web/applications/asr-rpc/internal"
	"github.com/mathiasXie/gin-web/applications/asr-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/consts"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	proto.UnimplementedASRServiceServer
	asrConfigs map[proto.Provider]asr.Config
}

func NewServer() *Server {

	asrConfigs := map[proto.Provider]asr.Config{
		proto.Provider_VOSK: {
			Model: config.Instance.ASR.Vosk.Model,
		},
		proto.Provider_ALIYUN: {
			APIKey: config.Instance.ASR.Aliyun.ApiKey,
		},
	}

	return &Server{
		asrConfigs: asrConfigs,
	}
}

func (s *Server) getASRProvider(ctx context.Context, provider proto.Provider) (asr.ASRProvider, context.Context, error) {
	asrConfig, ok := s.asrConfigs[provider]
	if !ok {
		logger.CtxError(ctx, "failed to create asr provider: %v", provider)
		return nil, ctx, fmt.Errorf("failed to create tts provider: %v", provider)
	}
	var asrProvider asr.ASRProvider
	var err error
	switch provider {
	case proto.Provider_VOSK:
		asrProvider, err = asr.NewVoskASR(asrConfig)
	}
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to create asr provider: %v", err)
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
	return asrProvider, ctx, nil
}

func (s *Server) SpeechToText(ctx context.Context, req *proto.SpeechToTextRequest) (*proto.SpeechToTextResponse, error) {

	asrProvider, ctx, err := s.getASRProvider(ctx, req.Provider)
	if err != nil {
		logger.CtxError(ctx, "failed to get asr provider: %v", err)
		return nil, fmt.Errorf("failed to get asr provider: %v", err)
	}

	text, err := asrProvider.SpeechToText(
		ctx,
		req.AudioData,
	)
	if err != nil {
		logger.CtxError(ctx, "failed to convert text to speech: %v", err)
		return nil, fmt.Errorf("failed to convert text to speech: %v", err)
	}
	logger.CtxInfo(ctx, "req: %v", req)
	return &proto.SpeechToTextResponse{
		Text: text,
	}, nil
}
