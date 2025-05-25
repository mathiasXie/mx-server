package handler

import (
	"fmt"

	llm_proto "github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

func (h *ChatHandler) IndentHandler() (text string, needLLM bool, afterFun func(), err error) {

	resp, err := (*h.LLMClient).IndentDetect(h.rpcCtx, &llm_proto.IndentRequest{
		Messages: h.userInfo.ChatMessages,
		Provider: llm_proto.LLMProvider(llm_proto.LLMProvider_value[config.Instance.Provider.Indent.LLM]),
		ModelId:  config.Instance.Provider.Indent.Model,
	})
	if err != nil {
		logger.CtxError(h.rpcCtx, "[ChatHandler]startToChat意图判断失败:", err)
		return "", false, nil, err
	}
	h.print(fmt.Sprintf("意图判断结果: %+v", resp), "blue")

	if resp.Nlu == "chat.continue" {
		return "", true, nil, nil
	}
	if resp.Nlu == "chat.weather" {
		return "上海今天天气晴朗，气温20度，风力3级，空气质量优,请使用这个数据回答用户", true, nil, nil
	}
	if resp.Nlu == "chat.exit" {

		return "再见", false, func() { h.conn.Close() }, nil
	}
	text = "哈哈哈哈"
	return text, false, nil, nil
}
