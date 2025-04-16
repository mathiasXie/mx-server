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

func InitFunctionRPC(rpcConfig *config.RPCConfig) FunctionRpcClient {
	if rpcConfig.Host == "" {
		return FunctionRpcClient{}
	}
	// 连接服务端
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", rpcConfig.Host, rpcConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := proto.NewFunctionServiceClient(conn)
	return FunctionRpcClient{
		FunctionServiceClient: client,
		Conn:                  conn,
	}
}
