package game

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"missevanbot/models"
	"missevanbot/utils"
)

type PassParcel struct {
	*Game
	Number int
}

func (p *PassParcel) Start(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	if p.State != models.StateReady {
		cmd.Output <- models.TplGameNotEnough
		return
	}

	p.State = models.StateRunning // state transfer

	go func() {
		// Sleep one minute for every five players.
		t := time.Minute * time.Duration(len(p.Players)/5+1)
		time.Sleep(t)
		// GAME OVER
		cmd.Output <- fmt.Sprintf(models.TplGameParcelOver, p.Players[p.Index].Name)
		// Add scores for winning players.
		if p.Index > 0 {
			addScore(cmd.ID, p.Players[:p.Index], models.ScorePassParcel)
		}
		if p.Index < len(p.Players)-1 {
			addScore(cmd.ID, p.Players[p.Index+1:], models.ScorePassParcel)
		}
		stop(cmd)
	}()

	text := strings.Builder{}
	text.WriteString(fmt.Sprintln(models.TplGameStart))
	text.WriteString(fmt.Sprintf(models.TplGameParcelNumber+"\n", p.parcelNumber()))
	text.WriteString(fmt.Sprintf(models.TplGameNext, p.Players[p.Index].Name))

	cmd.Output <- text.String()
}

func (p *PassParcel) Action(cmd *models.Command, textMsg models.FmTextMessage) {
	if p.State != models.StateRunning {
		return
	}

	if !p.isNowPlayer(textMsg.User.UserID) {
		return
	}

	if !p.isParcelCorrect(textMsg.Message) {
		cmd.Output <- models.TplGameInputIllegal
		return
	}
	p.nextPlayer()

	text := strings.Builder{}
	text.WriteString(fmt.Sprintf(models.TplGameParcelNumber+"\n", p.parcelNumber()))
	text.WriteString(fmt.Sprintf(models.TplGameNext, p.Players[p.Index].Name)) // output name of the next player.

	cmd.Output <- text.String()
}

// parcelNumber generates a random five-digit number,
// and stores it into m.
func (p *PassParcel) parcelNumber() string {
	number := utils.RandomNumber(10000, 99999)
	p.Number = number
	str := strconv.FormatInt(int64(number), 10)
	return str
}

// isParcelCorrect checks if int(s) equals the value stored in m.
func (p *PassParcel) isParcelCorrect(s string) bool {
	str := strconv.FormatInt(int64(p.Number), 10)
	return s == str
}
