package main

import (
	"missevan-fm/config"
	"missevan-fm/module"
)

// cronTask 定时任务
func cronTask(room *config.RoomConfig) {
	if room.Rainbow != nil && len(room.Rainbow) > 0 {
		go module.Praise(room)
	}
}
