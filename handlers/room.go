package handlers

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"missevanbot/config"
	"missevanbot/models"
	"missevanbot/modules"
)

func eventOpen(output chan<- string, room *models.Room) {
	data, err := models.NewTemplate(models.TmplStartUp, config.Nickname())
	if err != nil {
		zap.S().Warn(room.Log("create template failed", err))
	} else {
		output <- data
	}

	if room.Ticker == nil {
		room.Ticker = time.NewTicker(time.Minute)
	} else {
		room.Ticker.Reset(time.Minute)
	}
	room.TickerCount = 0

	if !room.Watch {
		return
	}
	text := fmt.Sprintf("%s 开播啦~", room.Creator)
	if err := modules.Push(modules.TitleOpen, text); err != nil {
		zap.S().Warn(room.Log("bark push failed", err))
	}
}

func eventClose(room *models.Room) {
	// stop the ticker.
	if t := room.Ticker; t != nil {
		t.Stop()
	}
	room.TickerCount = 0
	// clear playlist.
	room.Playlist = nil
	// stop the game instance.
	room.Gamer = nil
	// clear count.
	room.Count = 0

	if !room.Watch {
		return
	}
	// push notify to the message service.
	text := fmt.Sprintf("%s 下播啦~", room.Creator)
	if err := modules.Push(modules.TitleClose, text); err != nil {
		zap.S().Warn(room.Log("bark push failed", err))
	}
}
