package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

// InitRDBClient 初始化 Redis 客户端
func InitRDBClient(conf *RedisConfig) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Host,
		Password: conf.Password,
		DB:       conf.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(fmt.Errorf("连接Redis不成功: %s \n", err))
	}

	RDB = rdb
}
