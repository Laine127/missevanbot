package handler

import (
	"fmt"
	"time"

	"missevan-fm/config"
	"missevan-fm/module"
)

// handleRoom 处理直播间相关事件
func handleRoom(room *config.RoomConfig, textMsg *FmTextMessage) {
	conf := config.Conf.Push

	switch textMsg.Event {
	case EventStatistic:
		_statistics[room.ID].Online = textMsg.Statistics.Online
	case EventOpen:
		// 通知推送
		if r := module.RoomInfo(room.ID); r != nil {
			creatorName := r.Info.Creator.Username
			text := fmt.Sprintf("%s 开播啦~", creatorName)
			module.Push(conf, module.TitleOpen, text)
		}
		// 启用定时任务
		timer := time.NewTimer(1)
		_statistics[room.ID].Timer = timer
		go module.CronTask(room, timer)
	case EventClose:
		// 通知推送
		if r := module.RoomInfo(room.ID); r != nil {
			creatorName := r.Info.Creator.Username
			text := fmt.Sprintf("%s 下播啦~", creatorName)
			module.Push(conf, module.TitleClose, text)
		}
		// 关闭定时任务
		if timer := _statistics[room.ID].Timer; timer != nil {
			timer.Stop()
		}
	}
}
