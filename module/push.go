package module

import (
	"fmt"
	"net/http"

	"missevan-fm/bot"
)

const (
	TitleOpen  = "开播通知"
	TitleClose = "下播通知"
)

// Push 推送消息通知
func Push(title, msg string) {
	conf := bot.Conf.Push

	if conf.Bark != "" {
		if err := BarkPush(conf.Bark, title, msg); err != nil {
			ll.Print("Bark 推送失败", err.Error())
		}
	}
}

// BarkPush 通过 Bark 推送
func BarkPush(token, title, msg string) (err error) {
	_url := fmt.Sprintf("https://api.day.app/%s/%s/%s", token, title, msg)

	resp, err := http.Get(_url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return
}
