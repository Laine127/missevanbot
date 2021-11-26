package handler

import (
	"fmt"
	"log"
	"time"

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

// _cmdMap 帮助映射
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
func (room *RoomStore) handleCommand(cmdType int, textMsg *FmTextMessage) {
	roomID := room.Conf.ID
	switch cmdType {
	case CmdOnline:
		msg := fmt.Sprintf("当前直播间人数：%d~", room.Online)
		module.SendMessage(roomID, msg)
	case CmdSign:
		user := textMsg.User
		text, err := module.Sign(roomID, user.UserId, user.Username)
		if err != nil {
			log.Println("签到出错了。。。")
			return
		}
		msg := fmt.Sprintf("@%s %s", user.Username, text)
		module.SendMessage(roomID, msg)
	case CmdRank:
		if text := module.Rank(roomID); text != "" {
			msg := fmt.Sprintf("每日签到榜单：%s", text)
			module.SendMessage(roomID, msg)
		} else {
			module.SendMessage(roomID, "今天的榜单好像空空的~")
		}
	case CmdLove:
		module.SendMessage(roomID, "❤️~")
	case CmdBait:
		baitSwitch(room)
	case CmdHelper:
		fallthrough
	default:
		module.SendMessage(roomID, helpText)
	}
}

// baitSwitch 演员模式启停
func baitSwitch(room *RoomStore) {
	if room.Conf.Rainbow == nil || len(room.Conf.Rainbow) == 0 {
		// 设定的内容为空
		module.SendMessage(room.Conf.ID, "我好像没什么可说的。。。")
		return
	}
	if room.Bait && room.Timer != nil {
		module.SendMessage(room.Conf.ID, "我要去睡觉了")
		room.Bait = false
		room.Timer.Stop()
	} else {
		room.Bait = true
		if room.Conf.RainbowMaxInterval <= 0 {
			room.Conf.RainbowMaxInterval = 10
		}
		// 启用定时任务
		timer := time.NewTimer(1)
		room.Timer = timer
		go module.Praise(room.Conf, timer)
	}
}
