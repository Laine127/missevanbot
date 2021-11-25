package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mozillazg/go-pinyin"
	"missevan-fm/config"
	"missevan-fm/module"
)

type statistic struct {
	Count  int // 统计进入的数量
	Online int // 记录当前直播间在线人数
}

var _statistics = map[int]*statistic{}

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
	CmdLove
)

var _cmdMap = map[string]int{
	"帮助": CmdHelper,
	"在线": CmdOnline,
	"签到": CmdSign,
	"排行": CmdRank,
	"比心": CmdLove,
	"笔芯": CmdLove,
}

// handleTextMessage 处理文本消息
func handleTextMessage(room *config.RoomConfig, msg string) {
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

// handleRoom 处理直播间相关事件
func handleRoom(room *config.RoomConfig, textMsg *FmTextMessage) {
	conf := config.Conf.Push

	switch textMsg.Event {
	case EventStatistic:
		_statistics[room.ID].Online = textMsg.Statistics.Online
	case EventOpen:
		log.Println("直播间开启~")
		if conf.Bark != "" {
			// 通知推送
			if r := module.RoomInfo(room.ID); r != nil {
				creatorName := r.Info.Creator.Username
				text := fmt.Sprintf("%s 开播啦~", creatorName)
				module.Push(conf, module.TitleOpen, text)
			}
		}
	case EventClose:
		log.Println("直播间已经关闭了~")
		if conf.Bark != "" {
			// 通知推送
			if r := module.RoomInfo(room.ID); r != nil {
				creatorName := r.Info.Creator.Username
				text := fmt.Sprintf("%s 下播啦~", creatorName)
				module.Push(conf, module.TitleClose, text)
			}
		}
	}
}

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

// handleMessage 处理消息事件
func handleMessage(room *config.RoomConfig, textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventNew:
		if cmdType, ok := _cmdMap[textMsg.Message]; ok {
			// 命令处理
			handleCommand(room, cmdType, textMsg)
		}
	}
}

// handleCommand 处理消息中的命令
func handleCommand(room *config.RoomConfig, cmdType int, textMsg *FmTextMessage) {
	switch cmdType {
	case CmdOnline:
		module.SendMessage(room.ID, fmt.Sprintf("当前直播间人数：%d~", _statistics[room.ID].Online))
	case CmdSign:
		user := textMsg.User
		text, err := module.Sign(room.ID, user.UserId, user.Username)
		if err != nil {
			log.Println("签到出错了。。。")
			return
		}
		module.SendMessage(room.ID, fmt.Sprintf("@%s %s", user.Username, text))
	case CmdRank:
		if text := module.Rank(room.ID); text != "" {
			module.SendMessage(room.ID, fmt.Sprintf("每日签到榜单：%s", text))
		} else {
			module.SendMessage(room.ID, "今天的榜单好像空空的~")
		}
	case CmdLove:
		module.SendMessage(room.ID, "❤️~")
	case CmdHelper:
		fallthrough
	default:
		// 显示帮助提示
		module.SendMessage(room.ID, helperText)
	}
}
