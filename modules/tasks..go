package modules

import (
	"missevanbot/models"
	"missevanbot/modules/thirdparty"
)

const (
	DurPraise = 5
)

func RunTasks(output chan<- string, room *models.Room) {
	modes := ModeAll(room.ID)
	count := room.TickerCount

	if isEnabled(modes[ModeBait]) && shouldExec(count, DurPraise) {
		taskPraise(output)
	}
}

func taskPraise(output chan<- string) {
	if text, err := thirdparty.PraiseText(); err == nil && text != "" {
		output <- text
	}
}

func shouldExec(count int, dur int) bool {
	return count > 0 && count%dur == 0
}
