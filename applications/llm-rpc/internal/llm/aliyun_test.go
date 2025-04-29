package llm

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
)

func TestAliyunLLM_ChatStream(t *testing.T) {
	// 创建测试配置
	config := Config{
		APIKey:         "sk-69734b46cb6a4a24a0f29ab84f7451ae",
		BaseURL:        "https://dashscope.aliyuncs.com/compatible-mode/v1/",
		DefaultModelID: "qwen-plus",
	}

	// 创建LLM实例
	llm, err := NewAliyunLLM(config)
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
			Content: "你好",
		},
	}

	// 创建响应通道
	respChan := make(chan ChatStreamResponse)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 启动流式对话
	go func() {
		err := llm.ChatStream(ctx, config.DefaultModelID, messages, respChan)
		if err != nil {
			t.Errorf("流式对话失败: %v", err)
		}
	}()

	// 收集响应
	for {
		select {
		case resp := <-respChan:
			fmt.Printf("收到响应: %s\n", resp.Content)
			if resp.IsEnd {
				return
			}
		case <-ctx.Done():
			t.Fatal("测试超时")
		}
	}
}

func TestAliyunLLM_ChatNoStream(t *testing.T) {
	// 创建测试配置
	config := Config{
		APIKey:         "sk-69734b46cb6a4a24a0f29ab84f7451ae",
		BaseURL:        "https://dashscope.aliyuncs.com/compatible-mode/v1/",
		DefaultModelID: "qwen-plus",
	}

	// 创建LLM实例
	llm, err := NewAliyunLLM(config)
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
			Content: "你好",
		},
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := llm.ChatNoStream(ctx, config.DefaultModelID, messages)
	if err != nil {
		log.Fatalf("ChatNoStream Err:%v", err)
	}
	fmt.Printf("收到响应: %s\n", resp)
}
