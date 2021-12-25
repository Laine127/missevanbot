package handlers

import (
	"fmt"

	"go.uber.org/zap"
	"missevanbot/models"
	"missevanbot/modules"
	"missevanbot/utils"
)

func eventJoinQueue(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	for _, v := range textMsg.Queue {
		room.Count++
		if v.UserId == room.BotID() {
			continue
		}
		if name := v.Username; name != "" {
			var pinyin string
			if ok, err := modules.Mode(room.ID, modules.ModePinyin); err != nil {
				zap.S().Warn(room.Log("get mode failed", err))
			} else if ok {
				pinyin = utils.Pinyin(name)
			}

			data := struct {
				Name   string
				Pinyin string
				Extend string
			}{
				Name:   name,
				Pinyin: pinyin,
				Extend: modules.Word(modules.WordWelcome),
			}

			text, err := modules.NewTemplate(modules.TmplWelcome, data)
			if err != nil {
				zap.S().Warn(room.Log("create template failed", err))
				return
			}
			output <- text
		} else if room.Count > 5 && room.Count%2 == 0 {
			// Start sending welcome messages from the sixth person
			// and halve the number of messages sent.
			output <- models.TplWelcomeAnon
		}
	}
}

func eventFollowed(output chan<- string, textMsg models.FmTextMessage) {
	if name := textMsg.User.Username; name != "" {
		output <- fmt.Sprintf(models.TplThankFollow, name)
	}
}
