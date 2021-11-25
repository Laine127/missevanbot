package handler

import (
	"fmt"

	"missevan-fm/config"
	"missevan-fm/module"
)

// handleGift 处理礼物相关的事件
func handleGift(room *config.RoomConfig, textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventSend:
		// 有用户送礼物
		if username := textMsg.User.Username; username != "" {
			number := textMsg.Gift.Number
			giftName := textMsg.Gift.Name
			module.SendMessage(room.ID, fmt.Sprintf("感谢 @%s 赠送的%d个%s~", username, number, giftName))
		}
	}
}
