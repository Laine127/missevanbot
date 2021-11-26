package handler

import (
	"fmt"
	"log"
	"time"

	"missevan-fm/config"
	"missevan-fm/module"
)

// helpText 帮助文本
const helpText = `命令帮助：
帮助 -- 获取帮助信息
在线 -- 查看当前在线人数
签到 -- 在当前直播间进行签到
排行 -- 查看当前直播间当天签到排行`

const (
	CmdHelper = iota // 帮助提示
	CmdOnline        // 在线信息答复
	CmdSign          // 签到答复
	CmdRank          // 榜单答复
	CmdLove          // 比心答复
	CmdBait          // 演员模式启停
)

var _cmdMap = map[string]int{
	"帮助": CmdHelper,
	"在线": CmdOnline,
	"签到": CmdSign,
	"排行": CmdRank,
	// 下面是隐藏的命令
	"比心": CmdLove,
	"笔芯": CmdLove,
	"咳咳": CmdBait,
}

// handleCommand 处理消息中的命令
func handleCommand(room *config.RoomConfig, cmdType int, textMsg *FmTextMessage) {
	switch cmdType {
	case CmdOnline:
		msg := fmt.Sprintf("当前直播间人数：%d~", _statistics[room.ID].Online)
		module.SendMessage(room.ID, msg)
	case CmdSign:
		user := textMsg.User
		text, err := module.Sign(room.ID, user.UserId, user.Username)
		if err != nil {
			log.Println("签到出错了。。。")
			return
		}
		msg := fmt.Sprintf("@%s %s", user.Username, text)
		module.SendMessage(room.ID, msg)
	case CmdRank:
		if text := module.Rank(room.ID); text != "" {
			msg := fmt.Sprintf("每日签到榜单：%s", text)
			module.SendMessage(room.ID, msg)
		} else {
			module.SendMessage(room.ID, "今天的榜单好像空空的~")
		}
	case CmdLove:
		module.SendMessage(room.ID, "❤️~")
	case CmdBait:
		baitSwitch(room)
	case CmdHelper:
		fallthrough
	default:
		module.SendMessage(room.ID, helpText)
	}
}

// baitSwitch 演员模式启停
func baitSwitch(room *config.RoomConfig) {
	if room.Rainbow == nil || len(room.Rainbow) == 0 {
		module.SendMessage(room.ID, "我好像没什么可说的。。。")
		return
	}

	module.SendMessage(room.ID, "了解！")

	if s := _statistics[room.ID]; s.Bait && s.Timer != nil {
		s.Bait = false
		s.Timer.Stop()
	} else {
		s.Bait = true
		if room.RainbowMaxInterval <= 0 {
			room.RainbowMaxInterval = 10
		}

		// 启用定时任务
		timer := time.NewTimer(1)
		_statistics[room.ID].Timer = timer
		go module.Praise(room, timer)
	}
}
