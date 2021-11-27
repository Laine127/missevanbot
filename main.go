package main

import (
	"sync"

	"missevan-fm/bot"
)

func main() {
	bot.LoadConfig() // 读取并载入配置文件
	conf := bot.Conf
	bot.InitRDBClient(conf.Redis) // 初始化 Redis 缓存

	wg := &sync.WaitGroup{}

	for _, roomConf := range conf.Rooms {
		wg.Add(1)
		go func(conf *bot.RoomConfig) {
			connect(conf) // 主连接
		}(roomConf)
	}
	wg.Wait()
}
