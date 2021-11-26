package module

import (
	"math/rand"
	"time"

	"missevan-fm/config"
)

// Praise 彩虹屁模块
func Praise(room *config.RoomConfig, timer *time.Timer) {
	for {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randNumber := r.Intn(room.RainbowMaxInterval)
		timer.Reset(time.Duration(randNumber+1) * time.Minute)
		<-timer.C
		randIndex := r.Intn(len(room.Rainbow))
		SendMessage(room.ID, room.Rainbow[randIndex])
	}
}
