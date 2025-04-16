package loader

import (
	"context"
	"sync"

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
