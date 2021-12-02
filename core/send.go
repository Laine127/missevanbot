package core

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"missevan-fm/modules"
	"missevan-fm/utils"
)

type message struct {
	RoomID    int    `json:"room_id"`
	Message   string `json:"message"`
	MessageID string `json:"msg_id"`
}

// Send 从 outputMsg 中取出消息并发送
func Send(outputMsg <-chan string, roomID int) {
	for msg := range outputMsg {
		if msg == "" {
			continue
		}
		sendLoop(msg, roomID)
	}
}

func sendLoop(msg string, roomID int) {
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
