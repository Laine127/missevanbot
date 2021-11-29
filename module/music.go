package module

import (
	"fmt"
	"strconv"

	"missevan-fm/bot"
)

// MusicAdd 添加一首歌到歌单
func MusicAdd(roomID int, music string) {
	prefix := RedisPrefix + strconv.Itoa(roomID)
	bot.RDB.RPush(ctx, fmt.Sprintf("%s:musics", prefix), music)
}

// MusicAll 获取歌单中所有歌曲
func MusicAll(roomID int) []string {
	prefix := RedisPrefix + strconv.Itoa(roomID)
	musics := bot.RDB.LRange(ctx, fmt.Sprintf("%s:musics", prefix), 0, -1)
	return musics.Val()
}

// MusicPop 弹出一首歌曲
func MusicPop(roomID int) {
	prefix := RedisPrefix + strconv.Itoa(roomID)
	bot.RDB.LPop(ctx, fmt.Sprintf("%s:musics", prefix))
}

// MusicClear 清空歌单
func MusicClear(roomID int) {
	prefix := RedisPrefix + strconv.Itoa(roomID)
	bot.RDB.Del(ctx, fmt.Sprintf("%s:musics", prefix))
}
