package modules

import (
	"log"
	"testing"

	"missevan-fm/config"
)

func TestCookie(t *testing.T) {
	cookie, err := Cookie()
	log.Println(err)
	log.Println(cookie)
}

func TestFollow(t *testing.T) {
	config.LoadConfig()
	t.Run("test", func(t *testing.T) {
		ret, err := Follow(11111)
		log.Println(string(ret))
		log.Println(err)
	})
}
