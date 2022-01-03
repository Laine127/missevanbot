package handlers

import (
	"time"

	"go.uber.org/zap"
	"missevanbot/models"
	"missevanbot/modules"
)

const Delay = 3 * time.Second

func eventSend(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	user := textMsg.User
	if user.Username == "" {
		return
	}

	rid := textMsg.RoomID
	uid := user.UserID
	gift := textMsg.Gift
	// Check if gifts map exists, must execute before GiftAdd.
	exists := modules.GiftExists(rid, uid)
	modules.GiftAdd(rid, uid, gift.Name, gift.Number)
	if timer := room.GiftTimer; timer != nil {
		timer.Reset(Delay)
	}

	if exists {
		return // Do not create the new goroutine.
	}

	go func() {
		timer := time.NewTimer(Delay)
		room.GiftTimer = timer
		<-timer.C
		m := modules.GiftGetRemove(textMsg.RoomID, user.UserID)
		if m == nil || len(m) == 0 {
			return
		}
		data := struct {
			Sender string
			Gifts  map[string]string
		}{
			Sender: user.Username,
			Gifts:  m,
		}

		text, err := modules.NewTemplate(modules.TmplGift, data)
		if err != nil {
			zap.S().Warn(room.Log("create template failed", err))
			return
		}
		output <- text
	}()
}
