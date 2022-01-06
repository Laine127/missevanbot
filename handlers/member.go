package handlers

import (
	"fmt"

	"go.uber.org/zap"
	"missevanbot/models"
	"missevanbot/modules"
	"missevanbot/utils"
)

func eventJoinQueue(output chan<- string, room *models.Room, queue []models.FmQueue) {
	for _, v := range queue {
		room.Count++
		if v.UserID == room.BotID() {
			continue
		}
		if modules.IsGlobalInvisible(v.UserID) || modules.IsInvisible(room.ID, v.UserID) {
			continue // invisible online
		}
		if name := v.Username; name != "" {
			var pinyin, pinyinM string
			if mode(room, modules.ModePinyin) {
				pinyin, pinyinM = utils.Pinyin(name)
			}

			// Store the information of the joined user.
			var noble string
			var nobleLevel int
			var medal string // medal name
			// var medalLevel string
			var medalCreator string // username of the medal creator
			var level int
			for _, t := range v.Titles {
				switch t.Type {
				case models.TitleLevel:
					level = t.Level
				case models.TitleMedal:
					medal = t.Name
					// medalLevel = v.Level
				case models.TitleNoble:
					noble = t.Name
					nobleLevel = t.Level
				}
			}

			if mode(room, modules.ModeMedal) && medal != "" {
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
				PinyinM      string
				Extend       string
			}{
				Name:         name,
				Level:        level,
				Noble:        noble,
				NobleLevel:   nobleLevel,
				MedalCreator: medalCreator,
				Contribution: v.Contribution,
				Pinyin:       pinyin,
				PinyinM:      pinyinM,
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

func eventFollowed(output chan<- string, user models.FmUser) {
	if name := user.Username; name != "" {
		output <- fmt.Sprintf(models.TplThankFollow, name)
	}
}

func eventAddAdmin(output chan<- string, room *models.Room, target models.FmTarget) {
	text, err := modules.NewTemplate(modules.TmplNewAdmin, target.Username)
	if err != nil {
		zap.S().Warn(room.Log("create template failed", err))
		return
	}
	output <- text
}
