package core

import (
	"context"
	"time"

	"missevanbot/models"
	"missevanbot/modules"
)

func Cron(ctx context.Context, output chan<- string, room *models.Room) {
	room.Ticker = time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			room.TickerCount = (room.TickerCount + 1) % 60
			<-room.Ticker.C
			modules.RunTasks(output, room)
		}
	}
}
