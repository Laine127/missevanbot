package handler

import (
	"fmt"

	"missevan-fm/config"
	"missevan-fm/module"
)

// handleRoom 处理直播间相关事件
func handleRoom(room *config.RoomConfig, textMsg *FmTextMessage) {
	s := _statistics[room.ID]

	switch textMsg.Event {
	case EventStatistic:
		s.Online = textMsg.Statistics.Online
	case EventOpen:
		// 发送帮助信息
		module.SendMessage(room.ID, "芝士机器人在线了，可以在直播间输入“帮助”或者@我来获取支持哦～")
		// 通知推送
		if r := module.RoomInfo(room.ID); r != nil {
			creatorName := r.Info.Creator.Username
			text := fmt.Sprintf("%s 开播啦~", creatorName)
			module.Push(module.TitleOpen, text)
		}
	case EventClose:
		// 通知推送
		if r := module.RoomInfo(room.ID); r != nil {
			creatorName := r.Info.Creator.Username
			text := fmt.Sprintf("%s 下播啦~", creatorName)
			module.Push(module.TitleClose, text)
		}
		// 关闭定时任务
		s.Bait = false
		if timer := s.Timer; timer != nil {
			timer.Stop()
		}
	}
}
