package game

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"missevanbot/models"
	"missevanbot/modules"
)

type GuessWord struct {
	*Game
	Word string
}

func (g *GuessWord) Start(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	if g.State != models.StateReady {
		cmd.Output <- models.TplGameNotEnough
		return
	}

	g.State = models.StateRunning // state transfer

	g.setWord()

	text := strings.Builder{}
	text.WriteString(fmt.Sprintln(models.TplGameStart))
	text.WriteString(fmt.Sprintf(models.TplGameNext+"\n", g.Players[g.Index].Name))
	text.WriteString(models.TplGameGuessStart)

	uid := g.Players[g.Index].ID
	content := fmt.Sprintf(models.TplGameGuessWord, g.Word)
	if _, err := modules.SendMessage(uid, content); err != nil {
		zap.S().Warn(cmd.Room.Log("send private message failed", err))
		stop(cmd)
		cmd.Output <- models.TplSthWrong
		return
	}

	cmd.Output <- text.String()
}

func (g *GuessWord) Action(cmd *models.Command, textMsg models.FmTextMessage) {
	if g.State != models.StateRunning {
		return
	}

	if !g.isNowPlayer(textMsg.User.UserID) {
		return
	}

	text := strings.Builder{}

	if g.guess(textMsg.Message) && g.Index != 0 {
		cmd.Output <- fmt.Sprintf(models.TplGameGuessSuccess, g.Players[g.Index].Name)
		addScore(cmd.Room.ID, []models.Player{g.Players[g.Index]}, 15)
		stop(cmd)
		return
	} else if g.Index != 0 {
		text.WriteString(fmt.Sprintln(models.TplGameGuessFailed))
	}
	g.nextPlayer()
	text.WriteString(fmt.Sprintf(models.TplGameNext, g.Players[g.Index].Name))
	if g.Index == 0 {
		text.WriteString("\n" + models.TplGameGuessTip)
	}

	cmd.Output <- text.String()
}

func (g *GuessWord) setWord() {
	g.Word = modules.Word(modules.WordGuess)
}

func (g *GuessWord) guess(s string) bool {
	return s == g.Word
}
