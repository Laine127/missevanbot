package core

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"missevanbot/models"
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
func Send(ctx context.Context, output <-chan string, room *models.Room) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-output:
			if msg != "" {
				send(msg, room)
			}
		}
	}
}

func send(msg string, room *models.Room) {
	defer func() {
		if p := recover(); p != nil {
			zap.S().Error(room.Log("panic", p))
		}
	}()

	mute, err := modules.Mode(room.ID, modules.ModeMute)
	if err != nil {
		// Consider error as mute disabled mode.
		zap.S().Warn(room.Log("query mode mute failed", err))
	} else if mute {
		return
	}

	_url := "https://fm.missevan.com/api/chatroom/message/send"

	data, _ := json.Marshal(message{
		RoomID:    room.ID,
		Message:   msg,
		MessageID: utils.MessageID(),
	})

	header := http.Header{}
	header.Set("content-type", "application/json; charset=UTF-8")

	req := modules.NewRequest(_url, header, room.BotCookie, data)
	body, err := req.Post()
	if err != nil {
		zap.S().Warn(room.Log("send message failed", err))
		return
	}
	str := string(body)
	if str == `{"code":500150022,"info":"聊天内容含有违规信息"}` {
		send("检测到屏蔽词，消息被屏蔽了哦～", room)
		return
	}
	zap.S().Debug(room.Log(str, nil))
}
