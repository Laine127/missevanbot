package main

import (
	"encoding/json"
	"fmt"
	"log"

	"missevan-fm/config"
	"missevan-fm/module"
)

var count int  // 统计进入的数量
var online int // 记录当前直播间在线人数

// helperText 帮助文本
const helperText = `命令大全：
帮助 -- 获取帮助信息
在线 -- 查看当前在线人数
签到 -- 在当前直播间进行签到
排行 -- 查看当前直播间当天签到排行`

const (
	CmdHelper = iota
	CmdOnline
	CmdSign
	CmdRank
)

var _cmdMap = map[string]int{
	"帮助": CmdHelper,
	"在线": CmdOnline,
	"签到": CmdSign,
	"排行": CmdRank,
}

// handleTextMessage 处理文本消息
func handleTextMessage(roomID int, msg string) {
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

	switch textMsg.Type {
	case TypeNotify:
		// 过滤掉通知消息
		return
	case TypeMember:
		handleMember(roomID, textMsg)
	case TypeGift:
		handleGift(roomID, textMsg)
	case TypeMessage:
		handleMessage(roomID, textMsg)
	case TypeRoom:
		handleRoom(roomID, textMsg)
	default:
	}
	// fmt.Println(msg + "\n")
}

// handleRoom 处理直播间相关事件
func handleRoom(roomID int, textMsg *FmTextMessage) {
	conf := config.Conf.Push

	switch textMsg.Event {
	case EventStatistic:
		online = textMsg.Statistics.Online
	case EventOpen:
		log.Println("直播间开启~")
		if conf.Bark != "" {
			// 通知推送
			room := module.RoomInfo(roomID)
			if room == nil {
				return
			}
			creatorName := room.Info.Creator.Username
			text := fmt.Sprintf("%s 开播啦~", creatorName)
			module.Push(conf, module.TitleOpen, text)
		}
	case EventClose:
		log.Println("直播间已经关闭了~")
		if conf.Bark != "" {
			// 通知推送
			room := module.RoomInfo(roomID)
			if room == nil {
				return
			}
			creatorName := room.Info.Creator.Username
			text := fmt.Sprintf("%s 下播啦~", creatorName)
			module.Push(conf, module.TitleClose, text)
		}
	}
}

// handleMember 处理用户相关的事件
func handleMember(roomID int, textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventJoinQueue:
		// 有用户进入直播间
		for _, v := range textMsg.Queue {
			count++
			if count <= 1 {
				// 屏蔽第一个进入的匿名用户（可能是机器人自己）
				return
			}
			if username := v.Username; username != "" {
				module.SendMessage(roomID, fmt.Sprintf("欢迎 @%s 进入直播间~", username))
			} else {
				module.SendMessage(roomID, "欢迎新同学进入直播间~")
			}
		}
	case EventFollowed:
		// 有新关注
		if username := textMsg.User.Username; username != "" {
			module.SendMessage(roomID, fmt.Sprintf("感谢 @%s 的关注~", username))
		}
	}
}

// handleGift 处理礼物相关的事件
func handleGift(roomID int, textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventSend:
		// 有用户送礼物
		if username := textMsg.User.Username; username != "" {
			number := textMsg.Gift.Number
			giftName := textMsg.Gift.Name
			module.SendMessage(roomID, fmt.Sprintf("感谢 @%s 赠送的%d个%s~", username, number, giftName))
		}
	}
}

// handleMessage 处理消息事件
func handleMessage(roomID int, textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventNew:
		if cmdType, ok := _cmdMap[textMsg.Message]; ok {
			// 命令处理
			handleCommand(roomID, cmdType, textMsg)
		}
	}
}

// handleCommand 处理消息中的命令
func handleCommand(roomID, cmdType int, textMsg *FmTextMessage) {
	switch cmdType {
	case CmdOnline:
		module.SendMessage(roomID, fmt.Sprintf("当前直播间人数：%d~", online))
	case CmdSign:
		user := textMsg.User
		text, err := module.Sign(roomID, user.UserId, user.Username)
		if err != nil {
			log.Println("签到出错了。。。")
			return
		}
		module.SendMessage(roomID, fmt.Sprintf("@%s %s", user.Username, text))
	case CmdRank:
		if text := module.Rank(roomID); text != "" {
			module.SendMessage(roomID, fmt.Sprintf("每日签到榜单：%s", text))
		} else {
			module.SendMessage(roomID, "今天的榜单好像空空的~")
		}
	case CmdHelper:
		fallthrough
	default:
		// 显示帮助提示
		module.SendMessage(roomID, helperText)
	}
}
