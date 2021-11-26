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
		s := _statistics[room.ID]
		for _, v := range textMsg.Queue {
			s.Count++
			if username := v.Username; username != "" {
				module.SendMessage(room.ID, fmt.Sprintf("欢迎 @%s 进入直播间~", username))
				if room.Pinyin {
					// 如果注音功能开启了，发送注音消息
					py := pinyin.NewArgs()
					py.Style = pinyin.Tone
					module.SendMessage(room.ID, fmt.Sprintf("注音：%v", pinyin.Pinyin(username, py)))
				}
			} else if s.Count > 1 && s.Count%2 == 0 {
				// 屏蔽第一次匿名用户欢迎，减半欢迎匿名用户次数
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
