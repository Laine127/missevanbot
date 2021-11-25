package handler

import (
	"fmt"

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
	case EventClose:
		// 通知推送
		if r := module.RoomInfo(room.ID); r != nil {
			creatorName := r.Info.Creator.Username
			text := fmt.Sprintf("%s 下播啦~", creatorName)
			module.Push(conf, module.TitleClose, text)
		}
	}
}
