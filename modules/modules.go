package modules

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"missevanbot/config"
)

const (
	RedisPrefix          = "missevan:"
	RedisPrefixCookies   = RedisPrefix + "cookies:"
	RedisPrefixWords     = RedisPrefix + "words:"
	RedisPrefixTemplates = RedisPrefix + "templates:"
	RedisPrefixRunning   = RedisPrefix + "running:"
	RedisPrefixOnline    = RedisPrefix + "online:"
)

var rdb *redis.Client

var ctx = context.Background()

func Init() {
	rdb = config.RDB
}

// prefixRoom return the string `missevan:[RoomID]:`
func prefixRoom(rid int) string {
	return RedisPrefix + strconv.Itoa(rid) + ":"
}
