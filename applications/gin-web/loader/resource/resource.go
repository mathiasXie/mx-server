package resource

import (
	"github.com/mathiasXie/gin-web/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetResource() *DalResource {

	read, write := InitDB(&config.Instance.Resource.Mysql)
	redisClient := InitRedis(&config.Instance.Resource.Redis)
	ttsRpcClient := InitTTSRPC(&config.Instance.Resource.TTSRPC)
	return &DalResource{
		RedisClient:  redisClient,
		WriteDB:      write,
		ReadDB:       read,
		TTSRpcClient: ttsRpcClient,
	}
}

type DalResource struct {
	RedisClient  *redis.Client
	WriteDB      *gorm.DB
	ReadDB       *gorm.DB
	TTSRpcClient TTSRpcClient
}
