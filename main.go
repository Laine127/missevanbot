package main

import (
	"log"
	"os"
	"strconv"

	"missevan-fm/cache"
	"missevan-fm/config"
)

func main() {
	config.ReadConfig()                    // 读取配置文件
	cache.InitRDBClient(config.Conf.Redis) // 初始化 Redis 缓存

	rid := os.Args[1] // 455984011
	id, err := strconv.Atoi(rid)
	if err != nil {
		log.Println("Room ID错误。。。")
		return
	}
	connect(id) // 主连接
}
