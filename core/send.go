package core

import (
	"encoding/json"

	"go.uber.org/zap"
	"missevan-fm/handlers"
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
		sendLoop(msg, roomID)
	}
}

func sendLoop(msg string, roomID int) {
	defer func() {
		if p := recover(); p != nil {
			zap.S().Fatal(p)
		}
	}()

	_url := "https://fm.missevan.com/api/chatroom/message/send"

	data, _ := json.Marshal(message{
		RoomID:    roomID,
		Message:   msg,
		MessageID: utils.MessageID(),
	})
	if err := handlers.PostRequest(_url, data); err != nil {
		zap.S().Error(err)
	}
}
