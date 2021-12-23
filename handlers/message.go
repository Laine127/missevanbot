package handlers

import (
	"fmt"
	"strings"

	"missevanbot/config"
	"missevanbot/models"
	"missevanbot/modules"
)

func eventNew(output chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	first := strings.Split(textMsg.Message, " ")[0]
	if first == "" {
		return
	}
	// Determine whether it is a chat request and handle it.
	if first == fmt.Sprintf("@%s", modules.BotName()) || first == config.Nickname() {
		handleChat(output, textMsg)
		return
	}
	// Determine whether it is a command and handle it.
	if cmdType := models.Cmd(first); cmdType >= 0 {
		handleCommand(output, room, cmdType, textMsg)
		return
	}
	// If message is a game-related command or room.Gamer is not null,
	// that means there is a game instance exists, then handle the message.
	if cmdType := models.CmdGame(first); cmdType >= 0 || room.Gamer != nil {
		handleGame(output, room, textMsg, cmdType)
	}
	// Search for keywords in other messages and handle them.
	handleKeyword(output, textMsg)
}
