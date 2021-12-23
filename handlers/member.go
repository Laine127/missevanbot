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
		if username := v.Username; username != "" {
			var pinyin string
			if ok, err := modules.Mode(store.ID, modules.ModePinyin); err != nil {
				zap.S().Warn(store.Log("get mode failed", err))
			} else if ok {
				pinyin = utils.Pinyin(username)
			}

			sText := struct {
				Name   string
				Pinyin string
				Extend string
			}{username, pinyin, models.WelcomeString()}

			text, err := models.NewTemplate(models.TmplWelcome, sText)
			if err != nil {
				zap.S().Warn(store.Log("create template failed", err))
				return
			}
			output <- text
		} else if store.Count > 1 && store.Count%2 == 0 {
			// start sending welcome message from the second joined user,
			// and halve the number of messages sent.
			output <- models.TplWelcomeAnon
		}
	}
}

func eventFollowed(output chan<- string, textMsg models.FmTextMessage) {
	if username := textMsg.User.Username; username != "" {
		output <- fmt.Sprintf(models.TplThankFollow, username)
	}
}
