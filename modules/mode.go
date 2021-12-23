package modules

import (
	"errors"
	"fmt"
	"strconv"

	"missevanbot/config"
)

const (
	ModeMute   = "mute"
	ModePinyin = "pinyin"
	ModeBait   = "bait"
)

const (
	Enabled  = "1"
	Disabled = "0"
)

func InitMode(rid int) {
	rdb := config.RDB
	prefix := config.RedisPrefix + strconv.Itoa(rid) // Redis namespace prefix, `missevan:[RoomID]`
	key := fmt.Sprintf("%s:mode", prefix)            // `missevan:[RoomID]:mode`

	rdb.HSetNX(ctx, key, ModeMute, false)
	rdb.HSetNX(ctx, key, ModePinyin, false)
	rdb.HSetNX(ctx, key, ModeBait, false)
}

func Mode(rid int, mode string) (bool, error) {
	rdb := config.RDB
	prefix := config.RedisPrefix + strconv.Itoa(rid) // Redis namespace prefix, `missevan:[RoomID]`
	key := fmt.Sprintf("%s:mode", prefix)            // `missevan:[RoomID]:mode`

	cmd, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return false, err
	}

	switch m := cmd[mode]; m {
	case Enabled:
		return true, nil
	case Disabled:
		return false, nil
	case "":
		return false, errors.New("redis: empty mode key " + mode)
	default:
		return false, errors.New("redis: invalid mode value " + m)
	}
}

func ModeAll(rid int) map[string]string {
	rdb := config.RDB
	prefix := config.RedisPrefix + strconv.Itoa(rid)
	key := fmt.Sprintf("%s:mode", prefix)

	return rdb.HGetAll(ctx, key).Val()
}

func SetMode(rid int, mode string, val bool) {
	rdb := config.RDB
	prefix := config.RedisPrefix + strconv.Itoa(rid)
	key := fmt.Sprintf("%s:mode", prefix)

	rdb.HSet(ctx, key, mode, val)
}

func SwitchMode(rid int, mode string) {
	rdb := config.RDB
	prefix := config.RedisPrefix + strconv.Itoa(rid)
	key := fmt.Sprintf("%s:mode", prefix)

	switch m := rdb.HGet(ctx, key, mode).Val(); m {
	case Enabled:
		rdb.HMSet(ctx, key, mode, Disabled)
	case Disabled:
		rdb.HMSet(ctx, key, mode, Enabled)
	}
}

func isEnabled(s string) bool {
	return s != Disabled
}
