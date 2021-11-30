package modules

import (
	"log"

	"missevan-fm/config"
	"missevan-fm/modules/thirdparty"
)

const (
	TitleOpen  = "开播通知"
	TitleClose = "下播通知"
)

// Push 推送消息通知
func Push(title, msg string) {
	conf := config.Config().Push

	if conf.Bark != "" {
		if err := thirdparty.BarkPush(conf.Bark, title, msg); err != nil {
			log.Println("Bark 推送失败", err.Error())
		}
	}
}
