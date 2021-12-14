package modules

import (
	"fmt"
	"strconv"

	"missevan-fm/config"
)

// SongAdd add a song to the playlist.
func SongAdd(roomID int, music string) {
	prefix := config.RedisPrefix + strconv.Itoa(roomID)
	config.RDB.RPush(ctx, fmt.Sprintf("%s:playlist", prefix), music)
}

// SongAll return all the songs in the playlist.
func SongAll(roomID int) []string {
	prefix := config.RedisPrefix + strconv.Itoa(roomID)
	musics := config.RDB.LRange(ctx, fmt.Sprintf("%s:playlist", prefix), 0, -1)
	return musics.Val()
}

// SongPop pop up a song from the playlist.
func SongPop(roomID int) {
	prefix := config.RedisPrefix + strconv.Itoa(roomID)
	config.RDB.LPop(ctx, fmt.Sprintf("%s:playlist", prefix))
}

// SongClear clear the playlist.
func SongClear(roomID int) {
	prefix := config.RedisPrefix + strconv.Itoa(roomID)
	config.RDB.Del(ctx, fmt.Sprintf("%s:playlist", prefix))
}
