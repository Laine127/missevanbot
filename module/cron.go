package module

import (
	"time"

	"missevan-fm/config"
)

// CronTask 定时任务
func CronTask(room *config.RoomConfig, timer *time.Timer) {
	if room.Rainbow != nil && len(room.Rainbow) > 0 {
		go Praise(room, timer) // 彩虹屁任务
	}
}
