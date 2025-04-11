package resource

import (
	"fmt"

	"github.com/mathiasXie/gin-web/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis(redisConfig *config.RedisConfig) *redis.Client {

	if redisConfig.Host == "" {
		return nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port), // Redis服务器的地址
		Password: redisConfig.Password,                                     // 密码，没有则留空
		DB:       redisConfig.DB,                                           // 使用默认DB
	})
	return client
}
