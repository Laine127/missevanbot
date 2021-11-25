package handler

import (
	"fmt"
	"math/rand"
	"time"

	"missevan-fm/config"
	"missevan-fm/module"
)

var _chatList = [...]string{
	"你好呀～",
	"QAQ",
	"有什么事情嘛～",
	"不好意思呀我暂时没办法理你～",
	"抱歉现在有点忙～",
	"虽然很突然，其实我根本不是机器人",
	"你知道嘛，我的智商有25哦",
}

// handleChat 处理聊天请求
func handleChat(room *config.RoomConfig, textMsg *FmTextMessage) {
	r := rand.New(rand.NewSource(int64(room.ID) + time.Now().UnixNano()))
	idx := r.Intn(len(_chatList))
	module.SendMessage(room.ID, fmt.Sprintf("@%s %s", textMsg.User.Username, _chatList[idx]))
}
