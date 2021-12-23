package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"missevanbot/models"
)

type NumberBomb struct {
	*Game
	Bomb int
	Min  int
	Max  int
}

func (s *NumberBomb) Start(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	if s.State != models.StateReady {
		cmd.Output <- models.TplGameNotEnough
		return
	}

	s.State = models.StateRunning // state transfer

	s.Min, s.Max = s.bomb(len(s.Players))

	text := strings.Builder{}
	text.WriteString(fmt.Sprintln(models.TplGameStart))
	text.WriteString(fmt.Sprintf(models.TplGameBombRange+"\n", s.Min, s.Max))
	text.WriteString(fmt.Sprintf(models.TplGameNext, s.Players[s.Index].Name))

	cmd.Output <- text.String()
}

func (s *NumberBomb) Action(cmd *models.Command, textMsg models.FmTextMessage) {
	if s.State != models.StateRunning {
		return
	}

	if !s.isNowPlayer(textMsg.User.UserID) {
		return
	}

	text := strings.Builder{}
	if ok, min, max := s.guess(textMsg.Message); ok {
		cmd.Output <- fmt.Sprintf(models.TplGameBomb, textMsg.User.Username)
		// Add scores for winning players.
		if s.Index > 0 {
			addScore(cmd.Room.ID, s.Players[:s.Index], models.ScoreNumberBomb)
		}
		if s.Index < len(s.Players)-1 {
			addScore(cmd.Room.ID, s.Players[s.Index+1:], models.ScoreNumberBomb)
		}
		stop(cmd) // Stop the game.
		return
	} else if min == -1 {
		cmd.Output <- models.TplGameInputIllegal
		return
	} else {
		text.WriteString(fmt.Sprintf(models.TplGameBombRange+"\n", min, max))
	}

	s.nextPlayer()

	text.WriteString(fmt.Sprintf(models.TplGameNext, s.Players[s.Index].Name)) // output name of the next player.

	cmd.Output <- text.String()
}

// The bomb sets the range according to number of players,
// generates a random number in the range,
// stores values into m and returns the range boundary.
func (s *NumberBomb) bomb(players int) (int, int) {
	rand.Seed(time.Now().UnixNano())
	max := players * 30
	bomb := rand.Intn(max) + 1

	s.Min = 1
	s.Max = max
	s.Bomb = bomb

	return 1, max
}

// The guess parses g to int type, returns false and -1 if there is an error (not a number).
// If int(g) is not equals s.Bomb, return false and the new range boundary.
func (s *NumberBomb) guess(g string) (bool, int, int) {
	min := s.Min
	max := s.Max
	bomb := s.Bomb

	n, err := strconv.Atoi(g)
	if err != nil || n < min || n > max {
		return false, -1, -1
	}

	if n == bomb {
		return true, 0, 0
	} else if n < bomb {
		s.Min = n + 1
		return false, n + 1, max
	} else {
		s.Max = n - 1
		return false, min, n - 1
	}
}
