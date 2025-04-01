package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/common-nighthawk/go-figure"
	"github.com/mathiasXie/gin-web/applications/tts-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/applications/tts-rpc/server"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath = flag.String("f", "../../conf", "the config path")

func main() {
	flag.Parse()
	configFile := fmt.Sprintf("%s/%s.yaml", *configPath, "tts-rpc")
	err := config.Instance.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	logger.InitLog()
	myFigure := figure.NewFigure(config.Instance.AppName, "", true)
	myFigure.Print()

	// 创建gRPC服务器
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Instance.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterTTSServiceServer(s, server.NewServer())
	reflection.Register(s)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
