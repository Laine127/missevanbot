package modules

import (
	"fmt"
	"strconv"

	"missevan-fm/config"
)

// MusicAdd add a song to the playlist.
func MusicAdd(roomID int, music string) {
	prefix := config.RedisPrefix + strconv.Itoa(roomID)
	config.RDB.RPush(ctx, fmt.Sprintf("%s:musics", prefix), music)
}

// MusicAll get all the songs in the playlist.
func MusicAll(roomID int) []string {
	prefix := config.RedisPrefix + strconv.Itoa(roomID)
	musics := config.RDB.LRange(ctx, fmt.Sprintf("%s:musics", prefix), 0, -1)
	return musics.Val()
}

// MusicPop pop up a song in the playlist.
func MusicPop(roomID int) {
	prefix := config.RedisPrefix + strconv.Itoa(roomID)
	config.RDB.LPop(ctx, fmt.Sprintf("%s:musics", prefix))
}

// MusicClear clear the playlist.
func MusicClear(roomID int) {
	prefix := config.RedisPrefix + strconv.Itoa(roomID)
	config.RDB.Del(ctx, fmt.Sprintf("%s:musics", prefix))
}
