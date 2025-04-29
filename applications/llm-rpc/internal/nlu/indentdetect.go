package nlu

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mathiasXie/gin-web/applications/llm-rpc/internal/llm"

	"github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"go.uber.org/zap"
)

type IndentResp struct {
	Nlu          string       `json:"nlu"`
	FunctionCall FunctionCall `json:"function_call"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

var indentSystemPromp = `
		你是一个意图识别助手。请分析用户的最后一句话，判断用户意图属于以下哪一类：
		处理步骤:
		1. 思考意图类型.
		2. 生成function_call格式,arguments是转义的json字符串
		3. 返回纯JSON

		返回格式示例：
		1. 播放音乐意图: {"nlu":"media.play","function_call": {"name": "play_music", "arguments": "{\"song_name\": \"音乐名称\"}"}}
		2. 结束对话意图: {"nlu":"chat.exit","function_call": {"name": "handle_exit_intent", "arguments": "{\"say_goodbye\": \"goodbye\"}"}}
		3. 天气意图: {"nlu":"chat.weather","function_call": {"name": "get_weather", "arguments": {}}}
		4. 获取当天日期时间: {"nlu":"chat.date","function_call": {"name": "get_time"}}
		5. 继续聊天意图: {"nlu":"chat.continue","function_call": {"name": "continue_chat"}}

		注意:
		- 播放音乐：无歌名时，song_name设为"random"
		- 如果没有明显的意图，应按照继续聊天意图处理
		- 只返回纯JSON，不要任何其他内容
`

func IndentDetect(ctx context.Context, llm llm.LLMProvider, req *proto.IndentRequest) (*proto.IndentResponse, error) {
	startTime := time.Now()

	//替换掉message里的prompt
	for index, message := range req.Messages {
		if message.Role == proto.ChatMessageRole_SYSTEM {
			req.Messages[index].Content = indentSystemPromp
		}
	}

	// 调用LLM
	llmResp, err := llm.ChatNoStream(ctx, req.ModelId, req.Messages)
	if err != nil {
		return nil, err
	}
	// 处理返回
	intentStr := strings.TrimSpace(llmResp)
	// 提取纯JSON
	re := regexp.MustCompile(`\{.*\}`)
	match := re.FindString(intentStr)
	if match != "" {
		intentStr = match
	}
	// 尝试解析
	var intentData IndentResp
	if err := json.Unmarshal([]byte(intentStr), &intentData); err != nil {
		logger.CtxError(ctx, "解析意图失败", zap.String("response", intentStr))
		// 默认返回继续聊天
		return &proto.IndentResponse{
			Nlu:          "",
			FunctionCall: nil,
		}, nil
	}
	logger.CtxInfo(ctx, fmt.Sprintf("意图识别结束，%+v耗时：%s", intentData, time.Since(startTime)))
	return &proto.IndentResponse{
		Nlu: intentData.Nlu,
		FunctionCall: &proto.FunctionCall{
			Name:      intentData.FunctionCall.Name,
			Arguments: intentData.FunctionCall.Arguments,
		},
	}, nil
}
