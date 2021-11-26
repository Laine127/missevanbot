package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mozillazg/go-pinyin"
	"missevan-fm/config"
	"missevan-fm/module"
)

type statistic struct {
	Count  int         // 统计进入的数量
	Online int         // 记录当前直播间在线人数
	Bait   bool        // 是否开启演员模式
	Timer  *time.Timer // 定时任务计时器
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

// handleRoom 处理直播间相关事件
func handleRoom(room *config.RoomConfig, textMsg *FmTextMessage) {
	s := _statistics[room.ID]

	switch textMsg.Event {
	case EventStatistic:
		s.Online = textMsg.Statistics.Online
	case EventOpen:
		// 发送帮助信息
		module.SendMessage(room.ID, "芝士机器人在线了，可以在直播间输入“帮助”或者@我来获取支持哦～")
		// 通知推送
		if r := module.RoomInfo(room.ID); r != nil {
			creatorName := r.Info.Creator.Username
			text := fmt.Sprintf("%s 开播啦~", creatorName)
			module.Push(module.TitleOpen, text)
		}
	case EventClose:
		// 通知推送
		if r := module.RoomInfo(room.ID); r != nil {
			creatorName := r.Info.Creator.Username
			text := fmt.Sprintf("%s 下播啦~", creatorName)
			module.Push(module.TitleClose, text)
		}
		// 关闭定时任务
		s.Bait = false
		if timer := s.Timer; timer != nil {
			timer.Stop()
		}
	}
}

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

// handleGift 处理礼物相关的事件
func handleGift(room *config.RoomConfig, textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventSend:
		// 有用户送礼物
		if username := textMsg.User.Username; username != "" {
			gift := textMsg.Gift
			msg := fmt.Sprintf("感谢 @%s 赠送的%d个%s~", username, gift.Number, gift.Name)
			module.SendMessage(room.ID, msg)
		}
	}
}

// handleMessage 处理消息事件
func handleMessage(room *config.RoomConfig, textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventNew:
		if cmdType, ok := _cmdMap[textMsg.Message]; ok {
			// 判断是否是命令，进行处理
			handleCommand(room, cmdType, textMsg)
			return
		}
		if strings.HasPrefix(textMsg.Message, fmt.Sprintf("@%s", config.Conf.Name)) {
			// 判断是否是聊天请求，进行处理
			handleChat(room, textMsg)
			return
		}
	}
}
