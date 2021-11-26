package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/mozillazg/go-pinyin"
	"missevan-fm/config"
	"missevan-fm/module"
)

var rooms sync.Map // 存储每个直播间的实时信息

type RoomStore struct {
	Conf   *config.Room // 当前房间的配置
	Count  int          // 统计进入的数量
	Online int          // 记录当前直播间在线人数
	Bait   bool         // 是否开启演员模式
	Timer  *time.Timer  // 定时任务计时器
}

// HandleTextMessage 处理文本消息
func HandleTextMessage(conf *config.Room, msg string) {
	if msg == "❤️" {
		// 过滤掉心跳消息
		return
	}
	// 解析 JSON 数据
	textMsg := new(FmTextMessage)
	if err := json.Unmarshal([]byte(msg), textMsg); err != nil {
		log.Println("解析失败了：", err)
		log.Println(msg)
		return
	}
	// 如果不存在就存储一个新的
	v, _ := rooms.LoadOrStore(conf.ID, &RoomStore{
		Conf: conf,
	})
	room, ok := v.(*RoomStore)
	if !ok {
		log.Println("获取 RoomStore 错误了。。。")
		return
	}

	switch textMsg.Type {
	case TypeNotify:
		// 过滤掉通知消息
		return
	case TypeRoom:
		room.handleRoom(textMsg)
	case TypeMember:
		room.handleMember(textMsg)
	case TypeGift:
		room.handleGift(textMsg)
	case TypeMessage:
		room.handleMessage(textMsg)
	default:
	}
	fmt.Println(msg + "\n")
}

// handleRoom 处理直播间相关事件
func (room *RoomStore) handleRoom(textMsg *FmTextMessage) {
	roomID := room.Conf.ID
	switch textMsg.Event {
	case EventStatistic:
		room.Online = textMsg.Statistics.Online
	case EventOpen:
		// 发送帮助信息
		module.SendMessage(roomID, "芝士机器人在线了，可以在直播间输入“帮助”或者@我来获取支持哦～")
		// 通知推送
		if info := module.RoomInfo(roomID); info != nil {
			creatorName := info.Creator.Username
			text := fmt.Sprintf("%s 开播啦~", creatorName)
			module.Push(module.TitleOpen, text)
		}
	case EventClose:
		// 通知推送
		if info := module.RoomInfo(roomID); info != nil {
			creatorName := info.Creator.Username
			text := fmt.Sprintf("%s 下播啦~", creatorName)
			module.Push(module.TitleClose, text)
		}
		// 关闭定时任务
		room.Bait = false
		if timer := room.Timer; timer != nil {
			timer.Stop()
		}
	}
}

// handleMember 处理用户相关的事件
func (room *RoomStore) handleMember(textMsg *FmTextMessage) {
	roomID := room.Conf.ID
	switch textMsg.Event {
	case EventJoinQueue:
		// 有用户进入直播间
		for _, v := range textMsg.Queue {
			room.Count++
			if username := v.Username; username != "" {
				module.SendMessage(roomID, fmt.Sprintf("欢迎 @%s 进入直播间~", username))
				if room.Conf.Pinyin {
					// 如果注音功能开启了，发送注音消息
					py := pinyin.NewArgs()
					py.Style = pinyin.Tone
					module.SendMessage(roomID, fmt.Sprintf("注音：%v", pinyin.Pinyin(username, py)))
				}
			} else if room.Count > 1 && room.Count%2 == 0 {
				// 屏蔽第一次匿名用户欢迎，减半欢迎匿名用户次数
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
func (room *RoomStore) handleGift(textMsg *FmTextMessage) {
	roomID := room.Conf.ID
	switch textMsg.Event {
	case EventSend:
		// 有用户送礼物
		if username := textMsg.User.Username; username != "" {
			gift := textMsg.Gift
			msg := fmt.Sprintf("感谢 @%s 赠送的%d个%s~", username, gift.Number, gift.Name)
			module.SendMessage(roomID, msg)
		}
	}
}

// handleMessage 处理消息事件
func (room *RoomStore) handleMessage(textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventNew:
		// 判断是否是命令，进行处理
		if cmdType, ok := _cmdMap[textMsg.Message]; ok {
			room.handleCommand(cmdType, textMsg)
			return
		}
		// 判断是否是沟通请求，进行处理
		if strings.HasPrefix(textMsg.Message, fmt.Sprintf("@%s", config.Conf.Name)) {
			room.handleChat(textMsg)
			return
		}
	}
}
