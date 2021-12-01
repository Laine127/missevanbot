package modules

import (
	"math/rand"
	"time"

	"missevan-fm/config"
	"missevan-fm/modules/thirdparty"
)

// Praise 彩虹屁模块
func Praise(outputMsg chan<- string, room *config.RoomConfig, timer *time.Timer) {
	for {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randNumber := r.Intn(room.RainbowMaxInterval)
		timer.Reset(time.Duration(randNumber+1) * time.Minute)
		<-timer.C

		if text, err := thirdparty.PraiseText(); err == nil && text != "" {
			outputMsg <- text
		}
	}
}