package resource

import (
	"fmt"
	"log"

	"github.com/mathiasXie/gin-web/applications/asr-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/config"
	"google.golang.org/grpc"
)

type ASRRpcClient struct {
	proto.ASRServiceClient
	Conn *grpc.ClientConn
}

func InitASRRPC(rpcConfig *config.RPCConfig) ASRRpcClient {
	if rpcConfig.Host == "" {
		return ASRRpcClient{}
	}
	// 连接服务端
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", rpcConfig.Host, rpcConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := proto.NewASRServiceClient(conn)
	// // 创建带有元数据的上下文
	// trace_id, _ := ctx.Value(consts.LogID).(string)
	// md := metadata.Pairs("trace_id", trace_id)
	// rpcCtx := metadata.NewOutgoingContext(ctx, md)
	return ASRRpcClient{
		ASRServiceClient: client,
		Conn:             conn,
	}
}
