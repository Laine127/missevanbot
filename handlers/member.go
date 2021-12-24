package handlers

import (
	"fmt"

	"go.uber.org/zap"
	"missevanbot/models"
	"missevanbot/modules"
	"missevanbot/utils"
)

func eventJoinQueue(output chan<- string, store *models.Room, textMsg models.FmTextMessage) {
	for _, v := range textMsg.Queue {
		store.Count++
		if v.UserId == modules.BotID() {
			continue
		}
		if name := v.Username; name != "" {
			var pinyin string
			if ok, err := modules.Mode(store.ID, modules.ModePinyin); err != nil {
				zap.S().Warn(store.Log("get mode failed", err))
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
				zap.S().Warn(store.Log("create template failed", err))
				return
			}
			output <- text
		} else if store.Count > 5 && store.Count%2 == 0 {
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
