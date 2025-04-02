package resource

import (
	"fmt"
	"log"

	"github.com/mathiasXie/gin-web/applications/function-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/config"
	"google.golang.org/grpc"
)

type FunctionRpcClient struct {
	proto.FunctionServiceClient
	Conn *grpc.ClientConn
}

// InitTTSRPC 初始化TTS RPC客户端 连接到tts-rpc服务 实例化proto.TTSServiceClient
// 返回proto.TTSServiceClient
// 如果ttsRpcConfig.Host为空，则返回nil
func InitFunctionRPC(functionRpcConfig *config.FunctionRPCConfig) FunctionRpcClient {
	if functionRpcConfig.Host == "" {
		return FunctionRpcClient{}
	}
	// 连接服务端
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", functionRpcConfig.Host, functionRpcConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := proto.NewFunctionServiceClient(conn)
	// // 创建带有元数据的上下文
	// trace_id, _ := ctx.Value(consts.LogID).(string)
	// md := metadata.Pairs("trace_id", trace_id)
	// rpcCtx := metadata.NewOutgoingContext(ctx, md)
	return FunctionRpcClient{
		FunctionServiceClient: client,
		Conn:                  conn,
	}
}
