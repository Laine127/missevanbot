package modules

import (
	"strconv"
	"time"

	"missevanbot/models"
)

// InitRoomInfo store some information of the room in Redis when bot is starting.
func InitRoomInfo(rid int, name string, bot string) {
	key := prefixRoom(rid) + "info"

	rdb.HSet(ctx, key, "alias", name)
	rdb.HSet(ctx, key, "bot", bot)
}

// StatusRunning store a live key with 30 seconds (equals to Websocket heart beat duration) expiration in Redis,
// the status Live means the bot running normally.
func StatusRunning(rid int) {
	key := RedisPrefixRunning + strconv.Itoa(rid) // `missevan:running:[RoomID]`

	rdb.SetEX(ctx, key, "", 30*time.Second)
}

// StatusOnline store an open key with one minute expiration in Redis,
// the status Open means the live room is opening.
func StatusOnline(room *models.Room) {
	key := RedisPrefixOnline + strconv.Itoa(room.ID) // `missevan:online:[RoomID]`

	var game string
	if gamer := room.Gamer; gamer != nil {
		game = gamer.GameName()
	}

	vals := map[string]interface{}{
		"online": room.Online,
		"count":  room.Count,
		"game":   game,
	}
	rdb.HSet(ctx, key, vals)
	rdb.Expire(ctx, key, time.Minute)
}
