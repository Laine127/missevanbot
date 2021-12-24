package modules

import (
	"missevanbot/config"
	"missevanbot/modules/thirdparty"
)

const (
	TitleOpen  = "开播通知"
	TitleClose = "下播通知"
	TitleError = "错误通知"
)

func MustPush(title, msg string) {
	_ = Push(title, msg)
}

// Push pushes a message via third-party push library,
// support bark push currently.
func Push(title, msg string) (err error) {
	conf := config.Push()

	if conf.Bark != "" {
		if err = thirdparty.BarkPush(conf.Bark, title, msg); err != nil {
			return
		}
	}
	return nil
}
