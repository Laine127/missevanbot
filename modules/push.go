package modules

import (
	"missevanbot/config"
	"missevanbot/modules/thirdparty"
)

const (
	TitleOpen  = "开播通知"
	TitleClose = "下播通知"
)

// Push 推送消息通知
func Push(title, msg string) (err error) {
	conf := config.Push()

	if conf.Bark != "" {
		if err = thirdparty.BarkPush(conf.Bark, title, msg); err != nil {
			return
		}
	}
	return nil
}
