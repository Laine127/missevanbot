package modules

import (
	"fmt"
	"time"

	"missevanbot/config"
)

func InvisibleSet(rid, uid int, expiration time.Duration) {
	rdb := config.RDB
	key := config.RedisPrefix + fmt.Sprintf("%d:invisible:%d", rid, uid)

	rdb.Set(ctx, key, "", expiration)
}

func IsInvisible(rid, uid int) bool {
	rdb := config.RDB
	key := config.RedisPrefix + fmt.Sprintf("%d:invisible:%d", rid, uid)

	return rdb.Exists(ctx, key).Val() == 1
}

func InvisibleCancel(rid, uid int) {
	rdb := config.RDB
	key := config.RedisPrefix + fmt.Sprintf("%d:invisible:%d", rid, uid)

	rdb.Del(ctx, key)
}
