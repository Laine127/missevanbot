package modules

import (
	"strconv"

	"missevanbot/config"
	"missevanbot/models"
	"missevanbot/modules/thirdparty"
)

func RunTasks(output chan<- string, room *models.Room) {
	modes := ModeAll(room.ID)
	count := room.TickerCount

	if mode := modes[ModePander]; isEnabled(mode) && shouldExec(count, mode) {
		taskPander(output)
	}
	if mode := modes[ModeWater]; isEnabled(mode) && shouldExec(count, mode) {
		taskWater(output)
	}
}

func taskPander(output chan<- string) {
	if text, err := thirdparty.PanderText(); err == nil && text != "" {
		output <- text
	}
}

func taskWater(output chan<- string) {
	rdb := config.RDB
	key := config.RedisPrefix + "words:water"
	c := rdb.SRandMember(ctx, key)

	if text := c.Val(); text != "" {
		output <- text
	}
}

func shouldExec(count int, dur string) bool {
	d, err := strconv.Atoi(dur)
	if err != nil {
		return false
	}
	return count > 0 && count%d == 0
}
