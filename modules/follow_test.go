package modules

import (
	"log"
	"testing"

	"missevan-fm/config"
)

func TestFollow(t *testing.T) {
	config.LoadConfig()
	t.Run("test", func(t *testing.T) {
		ret, err := Follow(11111)
		log.Println(string(ret))
		log.Println(err)
	})
}
