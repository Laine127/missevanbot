package message

import (
	"testing"

	"missevan-fm/bot"
)

func TestMustSend(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		bot.LoadConfig()
		MustSend(455984011, "测试")
	})
}
