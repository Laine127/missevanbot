package handlers

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"missevanbot/config"
	"missevanbot/handlers/game"
	"missevanbot/models"
	"missevanbot/modules"
)

// The HandleRoom handles the events that related to live room.
func HandleRoom(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventStatistic:
		room.Online = textMsg.Statistics.Online // update the number of online members.
	case models.EventOpen:
		eventOpen(output, room)
	case models.EventClose:
		eventClose(room)
	}
}

// The HandleMember handles the event related to member.
func HandleMember(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventJoinQueue:
		eventJoinQueue(output, room, textMsg)
	case models.EventFollowed:
		eventFollowed(output, textMsg)
	}
}

// The HandleGift handles the event related to gift.
func HandleGift(output chan<- string, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventSend:
		eventSend(output, textMsg)
	}
}

// The HandleMessage handles the event related to message.
func HandleMessage(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventNew:
		eventNew(output, room, textMsg)
	}
}

// The handleChat handles the chat requests.
func handleChat(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	if textMsg.Message == room.BotNic {
		output <- fmt.Sprintf("@%s %s", textMsg.User.Username, modules.Word(modules.WordReply))
		return
	}
	output <- Chat(textMsg.User.Username)
}

// The handleCommand handles the user commands,
// simple logics are handling in this function, others in command.go.
func handleCommand(output chan<- string, room *models.Room, cmdType int, textMsg models.FmTextMessage) {
	info, err := modules.RoomInfo(room)
	if err != nil {
		zap.S().Warn(room.Log("fetch the room information failed", err))
		return
	}

	user := textMsg.User // message sender
	args := strings.Fields(strings.TrimSpace(textMsg.Message))
	cmd := &models.Command{
		Args:   args[1:],
		Room:   room,
		Info:   info,
		User:   user,
		Role:   role(info, user.UserID), // get the role of the message sender.
		Output: output,
	}

	if f, ok := _cmdMap[cmdType]; ok {
		f(cmd)
	}
}

// handleGame handles the game-related message.
func handleGame(output chan<- string, room *models.Room, textMsg models.FmTextMessage, cmdType int) {
	info, err := modules.RoomInfo(room)
	if err != nil {
		zap.S().Warn(room.Log("fetch the room information failed", err))
		return
	}

	user := textMsg.User
	args := strings.Fields(strings.TrimSpace(textMsg.Message))
	cmd := &models.Command{
		Args:   args[1:],
		Room:   room,
		Info:   info,
		User:   user,
		Role:   role(info, user.UserID),
		Output: output,
	}

	// these cases handle the game creating operation,
	// other operations will be handling in the default case.
	switch cmdType {
	case models.CmdGameNumberBomb:
		if room.Gamer != nil {
			output <- models.TplGameExists
			return
		}
		room.Gamer = &game.NumberBomb{Game: new(game.Game)}
		room.Gamer.Create(cmd)
		return
	case models.CmdGamePassParcel:
		if room.Gamer != nil {
			output <- models.TplGameExists
			return
		}
		room.Gamer = &game.PassParcel{Game: new(game.Game)}
		room.Gamer.Create(cmd)
		return
	case models.CmdGameGuessWord:
		if room.Gamer != nil {
			output <- models.TplGameExists
			return
		}
		room.Gamer = &game.GuessWord{Game: new(game.Game)}
		room.Gamer.Create(cmd)
	default:
		gamer := room.Gamer
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
func handleKeyword(output chan<- string, textMsg models.FmTextMessage) {
	if strings.Contains(textMsg.Message, "emo") {
		keyEmotional(output, textMsg.User)
	}
}

// role returns role of the user according to UID.
func role(info models.FmInfo, uid int) int {
	switch uid {
	case config.Admin():
		return models.RoleSuper
	case info.Creator.UserID:
		return models.RoleCreator
	default:
		// Determine whether it is a room administrator.
		for _, v := range info.Room.Members.Admin {
			if v.UserID == uid {
				return models.RoleAdmin
			}
		}
		return models.RoleMember
	}
}

func mode(room *models.Room, m string) (ok bool) {
	ok, err := modules.Mode(room.ID, m)
	if err != nil {
		zap.S().Warn(room.Log("get mode failed", err))
	}
	return
}
