package main

import (
	"log"
	"os"
	"sync"

	"missevan-fm/bot"
	"missevan-fm/handler"
	"missevan-fm/module"
	"missevan-fm/util"
)

func init() {
	_ = os.Mkdir("logs", 0755)
	_ = os.Chmod("logs", 0755)
}

func main() {
	bot.LoadConfig() // 读取并载入配置文件
	conf := bot.Conf
	bot.InitRDBClient(conf.Redis) // 初始化 Redis 缓存

	wg := &sync.WaitGroup{}

	for _, roomConf := range conf.Rooms {
		wg.Add(1)
		go func(roomConf *bot.RoomConfig) {
			defer wg.Done()
			ll, err := util.NewLogger(conf.Level, roomConf.ID)
			if err != nil {
				log.Println(err)
				return
			}
			store := &handler.RoomStore{
				RoomConfig: roomConf,
				Logger:     ll,
			}
			module.MustFollow(module.RoomInfo(store.ID).Creator.UserID) // 关注当前主播
			connect(store)                                              // 主连接
		}(roomConf)
	}
	wg.Wait()
}
