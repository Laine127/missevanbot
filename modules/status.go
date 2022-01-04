package modules

import (
	"fmt"
	"time"

	"missevanbot/config"
	"missevanbot/models"
)

// InitRoomInfo store some information of the room in Redis when bot is starting.
func InitRoomInfo(rid int, name string, bot string) {
	rdb := config.RDB
	key := config.RedisPrefix + fmt.Sprintf("%d:info", rid) // missevan:[rid]:info

	rdb.HSet(ctx, key, "alias", name)
	rdb.HSet(ctx, key, "bot", bot)
}

// StatusRunning store a live key with 30 seconds (equals to Websocket heart beat duration) expiration in Redis,
// the status Live means the bot running normally.
func StatusRunning(rid int) {
	rdb := config.RDB
	key := config.RedisPrefix + fmt.Sprintf("running:%d", rid) // missevan:live:[rid]

	rdb.SetEX(ctx, key, "", 30*time.Second)
}

// StatusOnline store an open key with one minute expiration in Redis,
// the status Open means the live room is opening.
func StatusOnline(room *models.Room) {
	rdb := config.RDB
	key := config.RedisPrefix + fmt.Sprintf("online:%d", room.ID) // missevan:open:[rid]

	var game string
	if gamer := room.Gamer; gamer != nil {
		game = gamer.GameName()
	}

	vals := map[string]interface{}{
		"count": room.Online,
		"game":  game,
	}
	rdb.HSet(ctx, key, vals)
	rdb.Expire(ctx, key, time.Minute)
}
