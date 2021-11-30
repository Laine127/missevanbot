package core

import (
	"go.uber.org/zap"
	"missevan-fm/handlers"
	"missevan-fm/models"
)

// Match 处理输出的消息
func Match(inputMsg <-chan models.FmTextMessage, outputMsg chan<- string, room *models.Room) {
	for msg := range inputMsg {
		matchLoop(msg, outputMsg, room)
	}
}

func matchLoop(msg models.FmTextMessage, outputMsg chan<- string, room *models.Room) {
	defer func() {
		if p := recover(); p != nil {
			zap.S().Fatal(p)
		}
	}()

	switch msg.Type {
	case models.TypeRoom:
		handlers.HandleRoom(outputMsg, room, msg)
	case models.TypeMember:
		handlers.HandleMember(outputMsg, room, msg)
	case models.TypeGift:
		handlers.HandleGift(outputMsg, room, msg)
	case models.TypeMessage:
		handlers.HandleMessage(outputMsg, room, msg)
	default:
	}
}
