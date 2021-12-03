package modules

import (
	"fmt"
	"strconv"

	"missevan-fm/config"
	"missevan-fm/models"
)

// MusicAdd 添加一首歌到歌单
func MusicAdd(roomID int, music string) {
	prefix := models.RedisPrefix + strconv.Itoa(roomID)
	config.RDB.RPush(ctx, fmt.Sprintf("%s:musics", prefix), music)
}

// MusicAll 获取歌单中所有歌曲
func MusicAll(roomID int) []string {
	prefix := models.RedisPrefix + strconv.Itoa(roomID)
	musics := config.RDB.LRange(ctx, fmt.Sprintf("%s:musics", prefix), 0, -1)
	return musics.Val()
}

// MusicPop 弹出一首歌曲
func MusicPop(roomID int) {
	prefix := models.RedisPrefix + strconv.Itoa(roomID)
	config.RDB.LPop(ctx, fmt.Sprintf("%s:musics", prefix))
}

// MusicClear 清空歌单
func MusicClear(roomID int) {
	prefix := models.RedisPrefix + strconv.Itoa(roomID)
	config.RDB.Del(ctx, fmt.Sprintf("%s:musics", prefix))
}
