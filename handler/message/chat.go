package message

import (
	"fmt"
	"math/rand"
	"time"
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

// Chat 简单回复
func Chat(username string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := r.Intn(len(_chatList))
	return fmt.Sprintf("@%s %s", username, _chatList[idx])
}
