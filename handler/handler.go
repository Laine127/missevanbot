package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"missevan-fm/config"
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
