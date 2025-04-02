package resource

import (
	"fmt"
	"log"

	"github.com/mathiasXie/gin-web/applications/tts-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/config"
	"google.golang.org/grpc"
)

type TTSRpcClient struct {
	proto.TTSServiceClient
	Conn *grpc.ClientConn
}

// InitTTSRPC 初始化TTS RPC客户端 连接到tts-rpc服务 实例化proto.TTSServiceClient
// 返回proto.TTSServiceClient
// 如果ttsRpcConfig.Host为空，则返回nil
func InitTTSRPC(ttsRpcConfig *config.TTSRPCConfig) TTSRpcClient {
	if ttsRpcConfig.Host == "" {
		return TTSRpcClient{}
	}
	// 连接服务端
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", ttsRpcConfig.Host, ttsRpcConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := proto.NewTTSServiceClient(conn)
	// // 创建带有元数据的上下文
	// trace_id, _ := ctx.Value(consts.LogID).(string)
	// md := metadata.Pairs("trace_id", trace_id)
	// rpcCtx := metadata.NewOutgoingContext(ctx, md)
	return TTSRpcClient{
		TTSServiceClient: client,
		Conn:             conn,
	}
}
