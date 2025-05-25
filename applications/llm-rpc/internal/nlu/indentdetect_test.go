package nlu

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/mathiasXie/gin-web/applications/llm-rpc/internal/llm"
	"github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

func TestIndentDetect(t *testing.T) {
	configFile := "../../conf/llm-rpc.yaml"

	err := config.Instance.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	logger.InitLog()

	// 创建测试配置
	config := llm.Config{
		APIKey:         "sk-69734b46cb6a4a24a0f29ab84f7451ae",
		BaseURL:        "https://dashscope.aliyuncs.com/compatible-mode/v1/",
		DefaultModelID: "qwen-plus",
	}

	// 创建LLM实例
	llm, err := llm.NewAliyunLLM(config)
	if err != nil {
		t.Fatalf("创建LLM实例失败: %v", err)
	}

	// 创建测试消息
	messages := []*proto.ChatMessage{
		{
			Role:    proto.ChatMessageRole_SYSTEM,
			Content: "你是一个来自古代的AI助手,请文言文回答用户的问题。",
		},
		{
			Role:    proto.ChatMessageRole_USER,
			Content: "帮我设置一个10点半的闹钟",
		},
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &proto.IndentRequest{
		ModelId:  "qwen-plus",
		Messages: messages,
	}

	resp, err := IndentDetect(ctx, llm, req)
	if err != nil {
		log.Fatalf("IndentDetect Err:%v", err)
	}
	fmt.Printf("返回内容：%+v\n", resp)

}
