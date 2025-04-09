package server

import (
	"context"
	"fmt"

	llm "github.com/mathiasXie/gin-web/applications/llm-rpc/internal"
	"github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/consts"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	proto.UnimplementedLLMServiceServer
	llmConfigs map[proto.LLMProvider]llm.Config
}

func NewServer() *Server {

	llmConfigs := map[proto.LLMProvider]llm.Config{
		proto.LLMProvider_VOLCENGINE: {
			APIKey:         config.Instance.LLM.VolcEngine.AuthToken,
			BaseURL:        config.Instance.LLM.VolcEngine.BaseURL,
			Models:         config.Instance.LLM.VolcEngine.Models,
			DefaultModelID: config.Instance.LLM.VolcEngine.DefaultModelID,
		},
		proto.LLMProvider_ALIYUN: {
			APIKey:         config.Instance.LLM.Aliyun.APIKey,
			BaseURL:        config.Instance.LLM.Aliyun.BaseURL,
			DefaultModelID: config.Instance.LLM.Aliyun.DefaultModelID,
			Models:         config.Instance.LLM.Aliyun.Models,
		},
	}

	return &Server{
		llmConfigs: llmConfigs,
	}
}

func (s *Server) getLLMProvider(ctx context.Context, provider proto.LLMProvider) (llm.LLMProvider, context.Context, error) {
	llmConfig, ok := s.llmConfigs[provider]
	if !ok {
		logger.CtxError(ctx, "failed to create llm provider: %v", provider)
		return nil, ctx, fmt.Errorf("failed to create llm provider: %v", provider)
	}
	var llmProvider llm.LLMProvider
	var err error
	switch provider {
	case proto.LLMProvider_VOLCENGINE:
		llmProvider, err = llm.NewVolcEngineLLM(llmConfig)
	case proto.LLMProvider_ALIYUN:
		llmProvider, err = llm.NewAliyunLLM(llmConfig)
	}
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to create llm provider: %v", err)
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
	return llmProvider, ctx, nil
}

func (s *Server) ChatStream(req *proto.ChatRequest, stream proto.LLMService_ChatStreamServer) error {
	llmProvider, ctx, err := s.getLLMProvider(stream.Context(), req.Provider)
	if err != nil {
		logger.CtxError(stream.Context(), "failed to get llm provider: %v", err)
		return fmt.Errorf("failed to get llm provider: %v", err)
	}
	logger.CtxInfo(stream.Context(), "llm provider: %v", llmProvider)
	logger.CtxInfo(stream.Context(), "got request: %v", req)

	// 创建响应通道
	respChan := make(chan llm.ChatStreamResponse)

	// 启动流式对话
	go func() {
		err := llmProvider.ChatStream(ctx, req.ModelId, req.Messages, respChan)
		if err != nil {
			logger.CtxError(stream.Context(), "failed to chat stream: %v", err)
			// 发送错误响应
			respChan <- llm.ChatStreamResponse{
				Content: fmt.Sprintf("Error: %v", err),
				IsEnd:   true,
			}
		}
	}()

	// 处理流式响应
	for {
		select {
		case resp := <-respChan:
			// 发送响应给客户端
			if err := stream.Send(&proto.ChatResponse{
				Content: resp.Content,
				IsEnd:   resp.IsEnd,
				Created: resp.Created,
				Model:   resp.Model,
				Id:      resp.ID,
			}); err != nil {
				logger.CtxError(stream.Context(), "failed to send response: %v", err)
				return err
			}

			// 如果是结束标记，返回
			if resp.IsEnd {
				return nil
			}
		case <-stream.Context().Done():
			// 客户端断开连接
			logger.CtxInfo(stream.Context(), "client disconnected")
			return nil
		}
	}
}

func (s *Server) ModelList(ctx context.Context, req *proto.ModelListRequest) (*proto.ModelListResponse, error) {

	llmProvider, ctx, err := s.getLLMProvider(ctx, req.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get llm provider: %v", err)
	}

	models, err := llmProvider.GetModelList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get models list: %v", err)
	}
	return &proto.ModelListResponse{Models: models}, nil
}
