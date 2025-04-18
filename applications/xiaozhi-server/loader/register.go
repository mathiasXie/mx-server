package loader

import (
	"context"
	"sync"

	asr_proto "github.com/mathiasXie/gin-web/applications/asr-rpc/proto/pb/proto"
	function_proto "github.com/mathiasXie/gin-web/applications/function-rpc/proto/pb/proto"
	llm_proto "github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"
	tts_proto "github.com/mathiasXie/gin-web/applications/tts-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/loader/resource"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var initCompleted bool

var initOnce sync.Once

var clients *resource.DalResource

func InitResource() error {
	initOnce.Do(func() {
		clients = resource.GetResource()
		initCompleted = true
	})
	return nil
}

func GetDB(ctx context.Context, isWrite ...bool) *gorm.DB {
	if !initCompleted {
		panic("resource not init")
	}
	if len(isWrite) > 0 && isWrite[0] {
		return clients.WriteDB.WithContext(ctx)
	}
	return clients.ReadDB.WithContext(ctx)
}

func GetRedis() *redis.Client {
	if !initCompleted {
		panic("resource not init")
	}
	return clients.RedisClient
}

func GetTTSRpc() *tts_proto.TTSServiceClient {
	if !initCompleted {
		panic("resource not init")
	}
	return &clients.TTSRpcClient.TTSServiceClient
}

func GetFunctionRpc() *function_proto.FunctionServiceClient {
	if !initCompleted {
		panic("resource not init")
	}
	return &clients.FunctionRpcClient.FunctionServiceClient
}

func GetLLMRpc() *llm_proto.LLMServiceClient {
	if !initCompleted {
		panic("resource not init")
	}
	return &clients.LLMRpcClient.LLMServiceClient
}

func GetASRRpc() *asr_proto.ASRServiceClient {
	if !initCompleted {
		panic("resource not init")
	}
	return &clients.ASRRpcClient.ASRServiceClient
}
