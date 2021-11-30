package main

import (
	"os"
	"sync"

	"missevan-fm/config"
	"missevan-fm/core"
	"missevan-fm/models"
	"missevan-fm/utils/logger"
)

func init() {
	_ = os.Mkdir("logs", 0755)
	_ = os.Chmod("logs", 0755)
}

var wg = &sync.WaitGroup{}

func main() {
	config.LoadConfig()

	conf := config.Config()
	config.InitRDBClient(conf.Redis)

	_ = logger.Init(conf.Log)

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
