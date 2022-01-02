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
			if mode(room, modules.ModePinyin) {
				pinyin = utils.Pinyin(name)
			}

			// Store the information of the joined user.
			var noble string
			var nobleLevel int
			var medal string // medal name
			// var medalLevel string
			var medalCreator string // username of the medal creator
			var level int
			for _, v := range v.Titles {
				switch v.Type {
				case models.TitleLevel:
					level = v.Level
				case models.TitleMedal:
					medal = v.Name
					// medalLevel = v.Level
				case models.TitleNoble:
					noble = v.Name
					nobleLevel = v.Level
				}
			}

			if mode(room, modules.ModeMedal) {
				medalInfo, err := modules.MedalInfo(medal)
				if err != nil {
					zap.S().Warn(room.Log(fmt.Sprintf("get medal[%s] info failed", medal), err))
				} else {
					medalCreator = medalInfo.Medal.CreatorName
				}
			}

			if !mode(room, modules.ModeNoble) {
				// Do not show the noble of joined user.
				noble = ""
			}

			data := struct {
				Name         string
				Level        int
				Noble        string
				NobleLevel   int
				MedalCreator string
				Contribution int
				Pinyin       string
				Extend       string
			}{
				Name:         name,
				Level:        level,
				Noble:        noble,
				NobleLevel:   nobleLevel,
				MedalCreator: medalCreator,
				Contribution: v.Contribution,
				Pinyin:       pinyin,
				Extend:       modules.Word(modules.WordWelcome),
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
