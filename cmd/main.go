package main

import (
	"context"
	"log"
	"sync"

	"missevanbot/config"
	"missevanbot/core"
	"missevanbot/models"
	"missevanbot/modules"
	"missevanbot/utils/logger"
)

var wg = &sync.WaitGroup{}

func main() {
	config.LoadConfig()
	conf := config.Config()

	// init the Redis client.
	config.InitRDBClient(conf.Redis)

	// init the logger.
	if err := logger.Init(conf.Log, conf.Level); err != nil {
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

		wg.Add(4)

		input := make(chan models.FmTextMessage, 1)
		output := make(chan string, 1)
		room := models.NewRoom(roomConf)

		ctx, cancel := context.WithCancel(ctx)

		modules.InitMode(room.ID)

		go core.Connect(ctx, cancel, input, room)
		go core.Match(ctx, input, output, room)
		go core.Send(ctx, output, room)
		go core.Cron(ctx, output, room)
	}

	defer wg.Done()

	wg.Wait()
}
