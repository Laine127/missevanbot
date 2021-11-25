package module

import (
	"fmt"
	"log"
	"net/http"

	"missevan-fm/config"
)

const (
	TitleOpen  = "开播通知"
	TitleClose = "下播通知"
)

// Push 推送消息通知
func Push(conf *config.PushConfig, title, msg string) {
	BarkPush(conf.Bark, title, msg)
}

// BarkPush 通过 Bark 推送
func BarkPush(token, title, msg string) {
	_url := fmt.Sprintf("https://api.day.app/%s/%s/%s", token, title, msg)

	resp, err := http.Get(_url)
	if err != nil {
		log.Println("Bark 推送失败：", err)
		return
	}
	defer resp.Body.Close()
}
