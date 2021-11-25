package handler

import (
	"fmt"
	"log"

	"missevan-fm/config"
	"missevan-fm/module"
)

// handleCommand 处理消息中的命令
func handleCommand(room *config.RoomConfig, cmdType int, textMsg *FmTextMessage) {
	switch cmdType {
	case CmdOnline:
		// 在线信息答复
		module.SendMessage(room.ID, fmt.Sprintf("当前直播间人数：%d~", _statistics[room.ID].Online))
	case CmdSign:
		// 签到答复
		user := textMsg.User
		text, err := module.Sign(room.ID, user.UserId, user.Username)
		if err != nil {
			log.Println("签到出错了。。。")
			return
		}
		module.SendMessage(room.ID, fmt.Sprintf("@%s %s", user.Username, text))
	case CmdRank:
		// 榜单答复
		if text := module.Rank(room.ID); text != "" {
			module.SendMessage(room.ID, fmt.Sprintf("每日签到榜单：%s", text))
		} else {
			module.SendMessage(room.ID, "今天的榜单好像空空的~")
		}
	case CmdLove:
		// 比心答复
		module.SendMessage(room.ID, "❤️~")
	case CmdHelper:
		// 帮助提示答复
		fallthrough
	default:
		module.SendMessage(room.ID, helperText)
	}
}
