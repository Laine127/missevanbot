package modules

import (
	"strconv"
)

// GiftAdd add number to the value that mapped by the key gift.
func GiftAdd(rid, uid int, gift string, number int) {
	key := prefixRoom(rid) + "gift:" + strconv.Itoa(uid)

	old := rdb.HGet(ctx, key, gift).Val()
	n, _ := strconv.Atoi(old)
	rdb.HSet(ctx, key, gift, n+number)
}

// GiftExists return if gifts map exists.
func GiftExists(rid, uid int) bool {
	key := prefixRoom(rid) + "gift:" + strconv.Itoa(uid)

	return rdb.Exists(ctx, key).Val() != 0
}

// GiftGetRemove return the gifts map (gift_name:gift_number),
// then delete the map.
func GiftGetRemove(rid, uid int) (m map[string]string) {
	key := prefixRoom(rid) + "gift:" + strconv.Itoa(uid)

	m = rdb.HGetAll(ctx, key).Val()
	rdb.Del(ctx, key)
	return
}
