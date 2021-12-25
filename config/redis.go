package config

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	RedisPrefix          = "missevan:"
	RedisPrefixCookies   = RedisPrefix + "cookies:"
	RedisPrefixWords     = RedisPrefix + "words:"
	RedisPrefixTemplates = RedisPrefix + "templates:"
)

var RDB *redis.Client

// InitRDBClient initialize the Redis client.
func InitRDBClient(conf *RedisConfig) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Host,
		Password: conf.Password,
		DB:       conf.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(fmt.Errorf("connect to Redis server failed: %s", err))
	}

	RDB = rdb
}
