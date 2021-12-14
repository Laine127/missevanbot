package main

import (
	"context"
	"log"
	"sync"

	"missevan-fm/config"
	"missevan-fm/core"
	"missevan-fm/models"
	"missevan-fm/modules"
	"missevan-fm/utils/logger"
)

var wg = &sync.WaitGroup{}

func main() {
	// load the configurations.
	config.LoadConfig()
	conf := config.Config()

	// init the Redis client.
	config.InitRDBClient(conf.Redis)

	// init the logger.
	if err := logger.Init(conf.Level); err != nil {
		log.Println("init logger failed: ", err)
		return
	}

	// init the bot information.
	modules.InitBot()

	ctx := context.Background()
	for _, roomConf := range conf.Rooms {
		if !roomConf.Enable {
			continue
		}

		wg.Add(3)

		input := make(chan models.FmTextMessage, 1)
		output := make(chan string, 1)
		room := &models.Room{RoomConfig: roomConf}
		ctx, cancel := context.WithCancel(ctx)

		go core.Connect(ctx, cancel, input, room.ID)
		go core.Match(ctx, input, output, room)
		go core.Send(ctx, output, room.ID)
	}

	defer wg.Done()

	wg.Wait()
}
