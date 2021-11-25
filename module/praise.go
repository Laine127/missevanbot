package module

import (
	"math/rand"
	"time"

	"missevan-fm/config"
)

// Praise 彩虹屁模块
func Praise(room *config.RoomConfig, timer *time.Timer) {
	roomID := room.ID
	texts := room.Rainbow
	for {
		r := rand.New(rand.NewSource(int64(roomID) + time.Now().UnixNano()))
		randNumber := r.Intn(room.RainbowMaxInterval)
		timer.Reset(time.Duration(randNumber) * time.Minute)
		<-timer.C
		randIndex := r.Intn(len(texts))
		SendMessage(roomID, texts[randIndex])
	}
}
