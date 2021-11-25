package module

import (
	"math/rand"
	"time"

	"missevan-fm/config"
)

// Praise 彩虹屁模块
func Praise(room *config.RoomConfig) {
	roomID := room.ID
	texts := room.Rainbow
	for {
		randNumber := rand.Intn(room.RainbowMaxInterval)
		time.Sleep(time.Duration(randNumber) * time.Minute)
		randIndex := rand.Intn(len(texts))
		SendMessage(roomID, texts[randIndex])
	}
}
