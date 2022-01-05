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
	cfg := config.Config()

	// Initialize the Redis client.
	config.InitRDBClient(cfg.Redis)
	// Initialize the logger.
	if err := logger.Init(cfg.Log, cfg.Level); err != nil {
		log.Println("init logger failed: ", err)
		return
	}
	modules.Init()

	ctx := context.Background()
	for _, rc := range cfg.Rooms {
		if !rc.Enable {
			continue
		}
		room := models.NewRoom(rc)
		err := modules.InitBot(room)
		if err != nil {
			log.Printf("init bot %s failed: %s\n", rc.Alias, err)
			continue
		}
		modules.InitRoomInfo(room.ID, room.Alias, room.BotNic)
		modules.InitMode(room.ID)

		wg.Add(4)

		input := make(chan models.FmTextMessage, 1)
		output := make(chan string, 1)

		ctx, cancel := context.WithCancel(ctx)

		go core.Connect(ctx, cancel, input, room)
		go core.Match(ctx, input, output, room)
		go core.Send(ctx, output, room)
		go core.Cron(ctx, output, room)
	}

	defer wg.Done()

	wg.Wait()
}
