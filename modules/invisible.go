package modules

import (
	"strconv"
	"time"
)

func InvisibleSet(rid, uid int, expiration time.Duration) {
	key := prefixRoom(rid) + "invisible:" + strconv.Itoa(uid)

	rdb.Set(ctx, key, "", expiration)
}

func IsInvisible(rid, uid int) bool {
	key := prefixRoom(rid) + "invisible:" + strconv.Itoa(uid)

	return rdb.Exists(ctx, key).Val() == 1
}

func InvisibleCancel(rid, uid int) {
	key := prefixRoom(rid) + "invisible:" + strconv.Itoa(uid)

	rdb.Del(ctx, key)
}
