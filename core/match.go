package core

import (
	"go.uber.org/zap"
	"missevan-fm/handlers"
	"missevan-fm/models"
)

// Match receive the message from the channel inputMsg,
// handle the message event and send results into the channel outputMsg.
func Match(inputMsg <-chan models.FmTextMessage, outputMsg chan<- string, room *models.Room) {
	for msg := range inputMsg {
		matchLoop(msg, outputMsg, room)
	}
}

func matchLoop(msg models.FmTextMessage, outputMsg chan<- string, room *models.Room) {
	defer func() {
		if p := recover(); p != nil {
			zap.S().Error(p)
		}
	}()

	switch msg.Type {
	case models.TypeNotify:
		return
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

	zap.S().Info(msg)
}
