package main

import (
	"log"
	"sync"

	"missevan-fm/bot"
	"missevan-fm/handler"
	"missevan-fm/util"
)

func main() {
	bot.LoadConfig() // 读取并载入配置文件
	conf := bot.Conf
	bot.InitRDBClient(conf.Redis) // 初始化 Redis 缓存

	wg := &sync.WaitGroup{}

	for _, roomConf := range conf.Rooms {
		wg.Add(1)
		go func(roomConf *bot.RoomConfig) {
			ll, err := util.NewLogger(conf.Level, roomConf.ID)
			if err != nil {
				log.Println(err)
				return
			}
			store := &handler.RoomStore{
				RoomConfig: roomConf,
				Logger:     ll,
			}
			connect(store) // 主连接
		}(roomConf)
	}
	wg.Wait()
}
