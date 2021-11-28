package module

import (
	"log"
	"sync"

	"missevan-fm/util"
)

var once = sync.Once{}

var ll *util.Logger

func init() {
	once.Do(func() {
		l, err := util.NewLogger("debug", 0)
		if err != nil {
			log.Println("error initialize logger: ", err)
			return
		}
		ll = l
	})
}
