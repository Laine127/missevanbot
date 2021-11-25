package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"missevan-fm/config"
)

// helperText 帮助文本
const helperText = `命令帮助：
帮助 -- 获取帮助信息
在线 -- 查看当前在线人数
签到 -- 在当前直播间进行签到
排行 -- 查看当前直播间当天签到排行`

const (
	CmdHelper = iota
	CmdOnline
	CmdSign
	CmdRank
	CmdLove
)

var _cmdMap = map[string]int{
	"帮助": CmdHelper,
	"在线": CmdOnline,
	"签到": CmdSign,
	"排行": CmdRank,
	// 下面是隐藏的命令
	"比心": CmdLove,
	"笔芯": CmdLove,
}

type statistic struct {
	Count  int // 统计进入的数量
	Online int // 记录当前直播间在线人数
}

var _statistics = map[int]*statistic{} // 存储每个直播间的实时信息

// HandleTextMessage 处理文本消息
func HandleTextMessage(room *config.RoomConfig, msg string) {
	if msg == "❤️" {
		// 过滤掉心跳消息
		return
	}
	textMsg := new(FmTextMessage)
	if err := json.Unmarshal([]byte(msg), textMsg); err != nil {
		log.Println("解析失败了：", err)
		log.Println(msg)
		return
	}

	// 初始化全局 Map
	if _, ok := _statistics[room.ID]; !ok {
		_statistics[room.ID] = new(statistic)
	}

	switch textMsg.Type {
	case TypeNotify:
		// 过滤掉通知消息
		return
	case TypeRoom:
		handleRoom(room, textMsg)
	case TypeMember:
		handleMember(room, textMsg)
	case TypeGift:
		handleGift(room, textMsg)
	case TypeMessage:
		handleMessage(room, textMsg)
	default:
	}
	fmt.Println(msg + "\n")
}
