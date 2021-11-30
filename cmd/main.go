package main

import (
	"log"
	"sync"

	"missevan-fm/config"
	"missevan-fm/core"
	"missevan-fm/models"
	"missevan-fm/utils/logger"
)

var wg = &sync.WaitGroup{}

func main() {
	// load the configurations
	config.LoadConfig()
	conf := config.Config()

	// init the Redis client
	config.InitRDBClient(conf.Redis)

	// init the logger
	if err := logger.Init(conf.Log); err != nil {
		log.Println("init logger failed: ", err)
		return
	}

	for _, roomConf := range conf.Rooms {
		inputMsg := make(chan models.FmTextMessage, 1)
		outputMsg := make(chan string, 1)

		room := &models.Room{
			RoomConfig: roomConf,
		}

		wg.Add(3)

		go core.Connect(inputMsg, room.ID)
		go core.Match(inputMsg, outputMsg, room)
		go core.Send(outputMsg, room.ID)
	}

	defer wg.Done()

	wg.Wait()
}
