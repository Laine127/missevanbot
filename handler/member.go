package handler

import (
	"fmt"

	"github.com/mozillazg/go-pinyin"
	"missevan-fm/config"
	"missevan-fm/module"
)

// handleMember 处理用户相关的事件
func handleMember(room *config.RoomConfig, textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventJoinQueue:
		// 有用户进入直播间
		val := _statistics[room.ID]
		for _, v := range textMsg.Queue {
			val.Count++
			if username := v.Username; username != "" {
				module.SendMessage(room.ID, fmt.Sprintf("欢迎 @%s 进入直播间~", username))
				if room.Pinyin {
					py := pinyin.NewArgs()
					py.Style = pinyin.Tone
					module.SendMessage(room.ID, fmt.Sprintf("注音：%v", pinyin.Pinyin(username, py)))
				}
			} else if val.Count > 1 && val.Count%2 == 0 {
				// 减半欢迎新用户次数
				module.SendMessage(room.ID, "欢迎新同学进入直播间~")
			}
		}
	case EventFollowed:
		// 有新关注
		if username := textMsg.User.Username; username != "" {
			module.SendMessage(room.ID, fmt.Sprintf("感谢 @%s 的关注~", username))
		}
	}
}
