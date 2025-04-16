package resource

import (
	"fmt"
	"log"

	"github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/config"
	"google.golang.org/grpc"
)

type LLMRpcClient struct {
	proto.LLMServiceClient
	Conn *grpc.ClientConn
}

// InitLLMRPC 初始化LLM RPC客户端 连接到llm-rpc服务 实例化proto.LLMServiceClient
// 返回proto.LLMServiceClient
// 如果llmRpcConfig.Host为空，则返回nil
func InitLLMRPC(rpcConfig *config.RPCConfig) LLMRpcClient {
	if rpcConfig.Host == "" {
		return LLMRpcClient{}
	}
	// 连接服务端
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", rpcConfig.Host, rpcConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := proto.NewLLMServiceClient(conn)
	return LLMRpcClient{
		LLMServiceClient: client,
		Conn:             conn,
	}
}
