package resource

import (
	"fmt"
	"log"

	"github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/config"
	"google.golang.org/grpc"
)

type LLMRPCConfig struct {
	proto.LLMServiceClient
	Conn *grpc.ClientConn
}

// InitLLMRPC 初始化LLM RPC客户端 连接到llm-rpc服务 实例化proto.LLMServiceClient
// 返回proto.LLMServiceClient
// 如果llmRpcConfig.Host为空，则返回nil
func InitLLMRPC(llmRpcConfig *config.LLMRPCConfig) LLMRPCConfig {
	if llmRpcConfig.Host == "" {
		return LLMRPCConfig{}
	}
	// 连接服务端
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", llmRpcConfig.Host, llmRpcConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := proto.NewLLMServiceClient(conn)
	// // 创建带有元数据的上下文
	// trace_id, _ := ctx.Value(consts.LogID).(string)
	// md := metadata.Pairs("trace_id", trace_id)
	// rpcCtx := metadata.NewOutgoingContext(ctx, md)
	return LLMRPCConfig{
		LLMServiceClient: client,
		Conn:             conn,
	}
}
