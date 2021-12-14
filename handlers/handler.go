package handlers

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"missevan-fm/config"
	"missevan-fm/models"
	"missevan-fm/modules"
	"missevan-fm/modules/game"
	"missevan-fm/utils"
)

// HandleRoom handle the event related to room.
func HandleRoom(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	info, err := modules.RoomInfo(room.ID)
	if err != nil {
		zap.S().Error("fetch the room information failed: ", err)
		return
	}

	switch textMsg.Event {
	case models.EventStatistic:
		room.Online = textMsg.Statistics.Online // update the number of online members.
	case models.EventOpen:
		output <- models.TplBotStart
		if room.Watch {
			text := fmt.Sprintf("%s 开播啦~", info.Creator.Username)
			if err := modules.Push(modules.TitleOpen, text); err != nil {
				zap.S().Error("Bark push failed: ", err)
			}
		}
	case models.EventClose:
		if room.Watch {
			text := fmt.Sprintf("%s 下播啦~", info.Creator.Username)
			if err := modules.Push(modules.TitleClose, text); err != nil {
				zap.S().Error("Bark push failed: ", err)
			}
		}
		// stop the timer.
		room.Bait = false
		if timer := room.Timer; timer != nil {
			timer.Stop()
		}
		// clear playlist.
		modules.MusicClear(room.ID)
		// clear count.
		room.Count = 0
	}
}

// HandleMember handle the event related to member.
func HandleMember(output chan<- string, store *models.Room, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventJoinQueue:
		for _, v := range textMsg.Queue {
			store.Count++
			if username := v.Username; username != "" {
				text := fmt.Sprintf(models.TplWelcome, username)
				if store.Pinyin {
					text += fmt.Sprintf("\n注音：[ %s ]", utils.Pinyin(username))
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

// HandleGift handle the event related to gift.
func HandleGift(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventSend:
		if username := textMsg.User.Username; username != "" {
			gift := textMsg.Gift
			output <- fmt.Sprintf(models.TplThankGift, username, gift.Number, gift.Name)
		}
	}
}

// HandleMessage handle the event related to message.
func HandleMessage(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventNew:
		first := strings.Split(textMsg.Message, " ")[0]
		if first == "" {
			return
		}
		// determine whether it is a command and handle it.
		if cmdType := models.Command(first); cmdType >= 0 {
			handleCommand(output, room, cmdType, textMsg)
			return
		}
		// determine whether it is a chat request and handle it.
		if first == fmt.Sprintf("@%s", modules.Name()) {
			handleChat(output, room, textMsg)
			return
		}
		// determine whether it is a game-related message and handle it.
		if gameType := game.Command(first); gameType >= 0 || room.GameStore.Game != game.Null {
			handleGame(output, room, textMsg)
		}
		// search for keywords in other messages and handle them.
		handleKeyword(output, room, textMsg)
	}
}

// handleCommand handle the commands,
// handle simple logic in this function, handle other logic in command.go.
func handleCommand(output chan<- string, store *models.Room, cmdType int, textMsg models.FmTextMessage) {
	info, err := modules.RoomInfo(store.ID)
	if err != nil {
		zap.S().Error("fetch the room information failed: ", err)
		return
	}

	user := textMsg.User // message sender
	args := strings.Fields(strings.TrimSpace(textMsg.Message))
	cmd := &command{
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
		output <- models.HelpText
	default:
		_cmdMap[cmdType](cmd)
	}
}

// handleChat handle the chat requests.
func handleChat(output chan<- string, store *models.Room, textMsg models.FmTextMessage) {
	output <- Chat(textMsg.User.Username)
}

// handleKeyword handle the message which contains keyword.
func handleKeyword(output chan<- string, store *models.Room, textMsg models.FmTextMessage) {
	if strings.Contains(textMsg.Message, "emo") {
		keyEmotional(output, textMsg.User)
	}
}

// handleGame handle the game-related message.
func handleGame(output chan<- string, store *models.Room, textMsg models.FmTextMessage) {
	info, err := modules.RoomInfo(store.ID)
	if err != nil {
		zap.S().Error("fetch the room information failed: ", err)
		return
	}

	user := textMsg.User
	args := strings.Fields(strings.TrimSpace(textMsg.Message))
	cmd := &command{
		Args:   args[1:],
		Room:   store,
		Info:   info,
		User:   user,
		Role:   role(info, user.UserID),
		Output: output,
	}

	switch gameType := game.Command(args[0]); gameType {
	case game.CmdHelper:
		output <- game.HelpText
	case game.CmdNumberBomb, game.CmdPassParcel:
		gameCreate(cmd, gameType)
	case game.CmdJoin:
		gameJoin(cmd, textMsg)
	case game.CmdStart:
		gameStart(cmd)
	case game.CmdStop:
		gameStop(cmd)
	case game.CmdPlayers:
		gamePlayers(cmd)
	case game.CmdRank:
		gameRank(cmd)
	default:
		gameAction(cmd, textMsg)
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
