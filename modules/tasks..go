package modules

import (
	"math/rand"
	"time"

	"missevanbot/config"
	"missevanbot/modules/thirdparty"
)

// Praise praise module, use timer to send messages regularly.
func Praise(output chan<- string, room *config.RoomConfig, timer *time.Timer) {
	for {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randNumber := r.Intn(room.RainbowMaxInterval)
		timer.Reset(time.Duration(randNumber+1) * time.Minute)
		<-timer.C

		if text, err := thirdparty.PraiseText(); err == nil && text != "" {
			output <- text
		}
	}
}
