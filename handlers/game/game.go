package game

import (
	"fmt"
	"strings"

	"missevanbot/models"
	"missevanbot/modules"
)

type Game struct {
	Name    string
	State   int             // game state.
	Players []models.Player // store a list of players in order, contains ID and name of each player.
	Index   int             // the index of next player.
}

func (s *Game) GameName() string {
	return s.Name
}

func (s *Game) AddPlayer(player models.Player) bool {
	// check if the player has joined the game
	for _, v := range s.Players {
		if v.ID == player.ID {
			return false
		}
	}

	s.Players = append(s.Players, player)
	return true
}

func (s *Game) AllPlayers(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	if len(s.Players) == 0 {
		cmd.Output <- models.TplGameNoPlayer
		return
	}

	text := strings.Builder{}
	text.WriteString("玩家列表：")
	for k, v := range s.Players {
		text.WriteString(fmt.Sprintf("\n%d. %s", k+1, v.Name))
	}

	cmd.Output <- text.String()
}

func (s *Game) nextPlayer() {
	s.Index = (s.Index + 1) % len(s.Players) // switch to the next player.
}

// isNowPlayer return whether the UID of player that specified by Index equals uid.
func (s *Game) isNowPlayer(uid int) bool {
	return s.Players[s.Index].ID == uid
}

func (s *Game) Create(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	if s.State != models.StateNotCreated {
		cmd.Output <- models.TplGameExists
		return
	}

	s.State = models.StateCreated // state transfer
	s.Players = make([]models.Player, 0)

	cmd.Output <- models.TplGameCreate
}

func (s *Game) Join(cmd *models.Command, textMsg models.FmTextMessage) {
	if s.State != models.StateCreated && s.State != models.StateReady {
		return
	}

	user := textMsg.User

	if !s.AddPlayer(models.Player{ID: user.UserID, Name: user.Username}) {
		cmd.Output <- models.TplGameJoinDup
		return
	}

	// if the number of players is not less than 2,
	// switch to the ready state
	if len(s.Players) >= 2 {
		s.State = models.StateReady
	}

	cmd.Output <- fmt.Sprintf(models.TplGameJoin, user.Username)
}

func (s *Game) Stop(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	if cmd.Gamer == nil {
		cmd.Output <- models.TplGameNull
		return
	}

	stop(cmd)

	cmd.Output <- models.TplGameStop
}

// addScore add score for all the players.
func addScore(roomID int, players []models.Player, score int) {
	for _, p := range players {
		modules.UpdateScore(roomID, p.ID, score)
	}
}

func stop(cmd *models.Command) {
	cmd.Gamer = nil
}
