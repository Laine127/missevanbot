package main

import (
	"sync"

	"missevan-fm/cache"
	"missevan-fm/config"
)

func main() {
	config.Load() // 读取并载入配置文件
	conf := config.Conf
	cache.InitRDBClient(conf.Redis) // 初始化 Redis 缓存

	wg := &sync.WaitGroup{}

	for _, roomConf := range conf.Rooms {
		wg.Add(1)
		go func(conf *config.Room) {
			connect(conf) // 主连接
		}(roomConf)
	}
	wg.Wait()
}
