package main

import (
	"sync"

	"missevan-fm/cache"
	"missevan-fm/config"
)

func main() {
	config.ReadConfig() // 读取配置文件
	conf := config.Conf
	cache.InitRDBClient(conf.Redis) // 初始化 Redis 缓存

	wg := &sync.WaitGroup{}

	for _, room := range conf.Rooms {
		wg.Add(1)
		go func(room *config.RoomConfig) {
			connect(room) // 主连接
		}(room)
	}
	wg.Wait()
}
