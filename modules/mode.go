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
	ModePander = "pander"
	ModeWater  = "water"
)

const (
	Enabled  = "1"
	Disabled = "0"
)

const (
	DefaultPander = 6
	DefaultWater  = 12
)

var _defaults = map[string]int{
	ModePander: DefaultPander,
	ModeWater:  DefaultWater,
}

// InitMode initialize all the modes that storing in Redis and not exists,
// if you want to add some new modes, initialize them after define the constants.
func InitMode(rid int) {
	rdb := config.RDB
	prefix := config.RedisPrefix + strconv.Itoa(rid) // Redis namespace prefix, `missevan:[RoomID]`
	key := fmt.Sprintf("%s:mode", prefix)            // `missevan:[RoomID]:mode`

	rdb.HSetNX(ctx, key, ModeMute, Disabled)
	rdb.HSetNX(ctx, key, ModePinyin, Enabled)
	rdb.HSetNX(ctx, key, ModePander, Disabled)
	rdb.HSetNX(ctx, key, ModeWater, Disabled)
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
	case Disabled:
		return false, nil
	case "":
		return false, errors.New("redis: empty mode key " + mode)
	default:
		return true, nil
	}
}

func ModeAll(rid int) map[string]string {
	rdb := config.RDB
	prefix := config.RedisPrefix + strconv.Itoa(rid)
	key := fmt.Sprintf("%s:mode", prefix)

	return rdb.HGetAll(ctx, key).Val()
}

func SetMode(rid int, mode string, val interface{}) {
	rdb := config.RDB
	prefix := config.RedisPrefix + strconv.Itoa(rid)
	key := fmt.Sprintf("%s:mode", prefix)

	rdb.HSet(ctx, key, mode, val)
}

// SwitchMode change mode into Disabled if enabled,
// if mode disabled and has a default value,
// change it to default value.
func SwitchMode(rid int, mode string) {
	rdb := config.RDB
	prefix := config.RedisPrefix + strconv.Itoa(rid)
	key := fmt.Sprintf("%s:mode", prefix)

	d, ok := _defaults[mode]

	switch m := rdb.HGet(ctx, key, mode).Val(); m {
	case Disabled:
		if ok {
			rdb.HMSet(ctx, key, mode, d)
		} else {
			rdb.HMSet(ctx, key, mode, Enabled)
		}
	case Enabled:
		fallthrough
	default:
		rdb.HMSet(ctx, key, mode, Disabled)
	}
}

func Validate(str string) bool {
	n, err := strconv.Atoi(str)
	return err == nil && n >= 0
}

func isEnabled(s string) bool {
	return s != Disabled
}
