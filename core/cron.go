package core

import (
	"context"
	"time"

	"go.uber.org/zap"
	"missevanbot/models"
	"missevanbot/modules"
)

func Cron(ctx context.Context, output chan<- string, room *models.Room) {
	room.Ticker = time.NewTicker(time.Minute)
	if !isOpening(room) {
		room.Ticker.Stop()
	}

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

func isOpening(room *models.Room) bool {
	info, err := modules.RoomInfo(room)
	if err != nil {
		zap.S().Warn(room.Log("fetch the room information failed", err))
		return true
	}

	return info.Room.Status.Open != 0
}
