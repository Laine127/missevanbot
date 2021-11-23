package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"missevan-fm/module"
)

var online int

// handleTextMessage 处理文本消息
func handleTextMessage(roomID int, msg string) {
	fmt.Println(msg)
	if msg == "❤️" {
		// 过滤掉心跳消息
		return
	}
	textMsg := new(FmTextMessage)
	if err := json.Unmarshal([]byte(msg), textMsg); err != nil {
		log.Println("解析失败了：", err)
		return
	}

	switch textMsg.Type {
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
}

// handleRoom 处理直播间相关事件
func handleRoom(roomID int, textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventStatistic:
		online = textMsg.Statistics.Online
	}
}

// handleMember 处理用户相关的事件
func handleMember(roomID int, textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventJoin:
		// 有用户进入直播间
		for _, v := range textMsg.Queue {
			if username := v.Username; username != "" {
				module.SendMessage(roomID, fmt.Sprintf("欢迎 @%s 进入直播间~", username))
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
		if strings.HasPrefix(textMsg.Message, "call") {
			// 命令处理
			handleCommand(roomID, textMsg.Message)
		}
	}
}

// handleCommand 处理消息中的命令
func handleCommand(roomID int, msg string) {
	cmd := msg[5:]
	fmt.Println(cmd)
	switch cmd {
	case "人数":
		module.SendMessage(roomID, fmt.Sprintf("当前直播间人数：%d~", online))
	}
}
