package core

import (
	"context"

	"go.uber.org/zap"
	"missevanbot/handlers"
	"missevanbot/models"
)

// Match receive the message from the channel input,
// handle the message event and send results into the channel output.
func Match(ctx context.Context, input <-chan models.FmTextMessage, output chan<- string, room *models.Room) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-input:
			match(msg, output, room)
		}
	}
}

func match(msg models.FmTextMessage, output chan<- string, room *models.Room) {
	defer func() {
		if p := recover(); p != nil {
			zap.S().Error(room.Log("panic", p))
		}
	}()

	switch msg.Type {
	case models.TypeNotify:
		return
	case models.TypeRoom:
		handlers.HandleRoom(output, room, msg)
	case models.TypeMember:
		handlers.HandleMember(output, room, msg)
	case models.TypeGift:
		handlers.HandleGift(output, msg)
	case models.TypeMessage:
		handlers.HandleMessage(output, room, msg)
	default:
	}
}
