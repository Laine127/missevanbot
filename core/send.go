package core

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"missevanbot/modules"
	"missevanbot/utils"
)

type message struct {
	RoomID    int    `json:"room_id"`
	Message   string `json:"message"`
	MessageID string `json:"msg_id"`
}

// Send keep taking messages from the channel output,
// send messages to the live room according to roomID on MissEvan.
func Send(ctx context.Context, output <-chan string, roomID int) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-output:
			if msg != "" {
				send(msg, roomID)
			}
		}
	}
}

func send(msg string, roomID int) {
	defer func() {
		if p := recover(); p != nil {
			zap.S().Error(p)
		}
	}()

	_url := "https://fm.missevan.com/api/chatroom/message/send"

	data, _ := json.Marshal(message{
		RoomID:    roomID,
		Message:   msg,
		MessageID: utils.MessageID(),
	})

	header := http.Header{}
	header.Set("content-type", "application/json; charset=UTF-8")

	if body, err := modules.PostRequest(_url, header, data); err != nil {
		zap.S().Error(err)
	} else {
		zap.S().Debug(string(body))
	}
}
