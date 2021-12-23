package handlers

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"missevanbot/config"
	"missevanbot/handlers/game"
	"missevanbot/models"
	"missevanbot/modules"
	"missevanbot/utils"
)

// The HandleRoom handles the events that related to live room.
func HandleRoom(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	info, err := modules.RoomInfo(room.ID)
	if err != nil {
		zap.S().Warn(room.Log("fetch the room information failed", err))
		return
	}

	switch textMsg.Event {
	case models.EventStatistic:
		room.Online = textMsg.Statistics.Online // update the number of online members.
	case models.EventOpen:
		data, err := models.NewTemplate(models.TmplStartUp, config.Nickname())
		if err != nil {
			zap.S().Warn(room.Log("create template failed", err))
		} else {
			output <- data
		}

		room.Ticker.Reset(time.Minute)
		room.TickerCount = 0

		if !room.Watch {
			return
		}
		text := fmt.Sprintf("%s 开播啦~", info.Creator.Username)
		if err := modules.Push(modules.TitleOpen, text); err != nil {
			zap.S().Warn(room.Log("bark push failed", err))
		}
	case models.EventClose:
		// stop the ticker.
		if t := room.Ticker; t != nil {
			t.Stop()
		}
		room.TickerCount = 0
		// clear playlist.
		room.Playlist = nil
		// stop the game instance.
		room.Gamer = nil
		// clear count.
		room.Count = 0

		if !room.Watch {
			return
		}
		// push notify to the message service.
		text := fmt.Sprintf("%s 下播啦~", info.Creator.Username)
		if err := modules.Push(modules.TitleClose, text); err != nil {
			zap.S().Warn(room.Log("bark push failed", err))
		}
	}
}

// The HandleMember handles the event related to member.
func HandleMember(output chan<- string, store *models.Room, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventJoinQueue:
		for _, v := range textMsg.Queue {
			store.Count++
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
	case models.EventFollowed:
		if username := textMsg.User.Username; username != "" {
			output <- fmt.Sprintf(models.TplThankFollow, username)
		}
	}
}

// The HandleGift handles the event related to gift.
//
// TODO: the number of gifts may be incorrect.
func HandleGift(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventSend:
		if username := textMsg.User.Username; username != "" {
			gift := textMsg.Gift
			output <- fmt.Sprintf(models.TplThankGift, username, gift.Number, gift.Name)
		}
	}
}

// The HandleMessage handles the event related to message.
func HandleMessage(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventNew:
		first := strings.Split(textMsg.Message, " ")[0]
		if first == "" {
			return
		}
		// determine whether it is a chat request and handle it.
		if first == fmt.Sprintf("@%s", modules.BotName()) || first == config.Nickname() {
			handleChat(output, room, textMsg)
			return
		}
		// determine whether it is a command and handle it.
		if cmdType := models.Cmd(first); cmdType >= 0 {
			handleCommand(output, room, cmdType, textMsg)
			return
		}
		// if message is a game-related command or room.Gamer is not null,
		// that means there is a game instance exists, then handle the message.
		if cmdType := models.CmdGame(first); cmdType >= 0 || room.Gamer != nil {
			handleGame(output, room, textMsg, cmdType)
		}
		// search for keywords in other messages and handle them.
		handleKeyword(output, room, textMsg)
	}
}

// The handleChat handles the chat requests.
func handleChat(output chan<- string, store *models.Room, textMsg models.FmTextMessage) {
	if textMsg.Message == config.Nickname() {
		output <- fmt.Sprintf("@%s %s", textMsg.User.Username, models.ReplyString())
		return
	}
	output <- Chat(textMsg.User.Username)
}

// The handleCommand handles the user commands,
// simple logics are handling in this function, others in command.go.
func handleCommand(output chan<- string, store *models.Room, cmdType int, textMsg models.FmTextMessage) {
	info, err := modules.RoomInfo(store.ID)
	if err != nil {
		zap.S().Warn(store.Log("fetch the room information failed", err))
		return
	}

	user := textMsg.User // message sender
	args := strings.Fields(strings.TrimSpace(textMsg.Message))
	cmd := &models.Command{
		Args:   args[1:],
		Room:   store,
		Info:   info,
		User:   user,
		Role:   role(info, user.UserID), // get the role of the message sender.
		Output: output,
	}

	switch cmdType {
	case models.CmdLove:
		output <- "❤️~"
	case models.CmdHelper:
		text, err := models.NewTemplate(models.TmplHelper, nil)
		if err != nil {
			zap.S().Warn(cmd.Room.Log("create template failed", err))
			return
		}
		output <- text
	default:
		_cmdMap[cmdType](cmd)
	}
}

// handleGame handles the game-related message.
func handleGame(output chan<- string, store *models.Room, textMsg models.FmTextMessage, cmdType int) {
	info, err := modules.RoomInfo(store.ID)
	if err != nil {
		zap.S().Warn(store.Log("fetch the room information failed", err))
		return
	}

	user := textMsg.User
	args := strings.Fields(strings.TrimSpace(textMsg.Message))
	cmd := &models.Command{
		Args:   args[1:],
		Room:   store,
		Info:   info,
		User:   user,
		Role:   role(info, user.UserID),
		Output: output,
	}

	// these cases handle the game creating operation,
	// other operations will be handling in the default case.
	switch cmdType {
	case models.CmdGameNumberBomb:
		if store.Gamer != nil {
			output <- models.TplGameExists
			return
		}
		store.Gamer = &game.NumberBomb{Game: new(game.Game)}
		store.Gamer.Create(cmd)
		return
	case models.CmdGamePassParcel:
		if store.Gamer != nil {
			output <- models.TplGameExists
			return
		}
		store.Gamer = &game.PassParcel{Game: new(game.Game)}
		store.Gamer.Create(cmd)
		return
	case models.CmdGameGuessWord:
		if store.Gamer != nil {
			output <- models.TplGameExists
			return
		}
		store.Gamer = &game.GuessWord{Game: new(game.Game)}
		store.Gamer.Create(cmd)
	default:
		gamer := store.Gamer
		if gamer == nil {
			output <- models.TplGameNull
			return
		}
		switch cmdType {
		case models.CmdGameJoin:
			gamer.Join(cmd, textMsg)
		case models.CmdGameStart:
			gamer.Start(cmd)
		case models.CmdGameStop:
			gamer.Stop(cmd)
		case models.CmdGamePlayers:
			gamer.AllPlayers(cmd)
		default:
			gamer.Action(cmd, textMsg)
		}
	}
}

// handleKeyword handles the message which contains keyword.
func handleKeyword(output chan<- string, store *models.Room, textMsg models.FmTextMessage) {
	if strings.Contains(textMsg.Message, "emo") {
		keyEmotional(output, textMsg.User)
	}
}

// role return the role of the user which specified by userID.
func role(info models.FmInfo, userID int) int {
	switch userID {
	case config.Admin():
		// bot administrator.
		return models.RoleSuper
	case info.Creator.UserID:
		// room creator.
		return models.RoleCreator
	default:
		// determine whether it is an room administrator
		for _, v := range info.Room.Members.Admin {
			if v.UserID == userID {
				return models.RoleAdmin
			}
		}
		// general member.
		return models.RoleMember
	}
}
