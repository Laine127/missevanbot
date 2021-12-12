package handlers

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"missevan-fm/models"
	"missevan-fm/modules"
	"missevan-fm/modules/game"
)

// gameCreate initialize a GameStore instance,
// current state must be StateNotCreated.
func gameCreate(cmd *command, gameType int) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	gs := cmd.Room.GameStore
	if gs.State != game.StateNotCreated {
		cmd.Output <- models.TplGameExists
		return
	}

	switch gameType {
	case game.CmdNumberBomb:
		gs.Game = game.NumberBomb
	case game.CmdPassParcel:
		gs.Game = game.PassParcel
	}

	gs.State = game.StateCreated // state transfer
	gs.Players = make([]game.Player, 0)
	gs.Value = make(map[string]int)

	cmd.Output <- models.TplGameCreate
}

// gameJoin put message sender into the player list,
// current state must be StateCreated,
// if the number of players not less than 2 then switch to the state StateReady.
func gameJoin(cmd *command, textMsg models.FmTextMessage) {
	gs := cmd.Room.GameStore
	if gs.State != game.StateCreated && gs.State != game.StateReady {
		return
	}

	user := textMsg.User
	if !gs.AddPlayer(game.Player{ID: user.UserID, Name: user.Username}) {
		cmd.Output <- models.TplGameJoinDup
		return
	}

	// if the number of players is not less than 2,
	// switch to the ready state
	if len(gs.Players) >= 2 {
		gs.State = game.StateReady
	}

	cmd.Output <- fmt.Sprintf(models.TplGameJoin, user.Username)
}

// gameStart start the game, current state must be StateReady,
// if there is no error then switch to the state StateRunning.
func gameStart(cmd *command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	gs := cmd.Room.GameStore
	if gs.State != game.StateReady {
		if len(gs.Players) < 2 {
			cmd.Output <- models.TplGameNotEnough
		}
		return
	}
	gs.State = game.StateRunning // state transfer
	gs.Index = 0                 // set the first player.

	var min, max int // set for number bomb game.
	switch gs.Game {
	case game.NumberBomb:
		min, max = game.BombGenerate(gs.Value, len(gs.Players))
	case game.PassParcel:
		go func() {
			// set timer, add one minutes every five players.
			timer := time.NewTimer(time.Minute * time.Duration(len(gs.Players)/5+1))
			<-timer.C
			// game over.
			cmd.Output <- fmt.Sprintf(models.TplGameParcelOver, gs.Players[gs.Index].Name)
			// add score for winning players.
			if gs.Index > 0 {
				addScore(cmd.Room.ID, gs.Players[:gs.Index], game.ScorePassParcel)
			}
			if gs.Index < len(gs.Players)-1 {
				addScore(cmd.Room.ID, gs.Players[gs.Index+1:], game.ScorePassParcel)
			}
			// stop that game.
			stop(cmd.Room.GameStore)
		}()
	}

	text := strings.Builder{}
	text.WriteString(fmt.Sprintln(models.TplGameStart))

	switch gs.Game {
	case game.NumberBomb:
		text.WriteString(fmt.Sprintf(models.TplGameBombRange+"\n", min, max))
	case game.PassParcel:
		text.WriteString(fmt.Sprintf(models.TplGameParcelNumber+"\n", game.ParcelNumber(gs.Value)))
	}

	text.WriteString(fmt.Sprintf(models.TplGameNext, gs.Players[0].Name))

	cmd.Output <- text.String()
}

// gameAction handle all game actions, current state must be StateRunning,
// message sender must be the player pointed to by the index.
func gameAction(cmd *command, textMsg models.FmTextMessage) {
	gs := cmd.Room.GameStore
	if gs.State != game.StateRunning {
		return
	}

	if gs.Players[gs.Index].ID != textMsg.User.UserID {
		return
	}

	text := strings.Builder{}
	switch gs.Game {
	case game.NumberBomb:
		if ok, min, max := game.BombGuess(gs.Value, textMsg.Message); ok {
			cmd.Output <- fmt.Sprintf(models.TplGameBomb, textMsg.User.Username)
			// add score for winning players.
			if gs.Index > 0 {
				addScore(cmd.Room.ID, gs.Players[:gs.Index], game.ScoreNumberBomb)
			}
			if gs.Index < len(gs.Players)-1 {
				addScore(cmd.Room.ID, gs.Players[gs.Index+1:], game.ScoreNumberBomb)
			}
			// stop the game.
			stop(cmd.Room.GameStore)
			return
		} else if min == -1 {
			cmd.Output <- models.TplGameInputIllegal
			return
		} else {
			text.WriteString(fmt.Sprintf(models.TplGameBombRange+"\n", min, max))
		}
	case game.PassParcel:
		if !game.IsParcelCorrect(gs.Value, textMsg.Message) {
			cmd.Output <- models.TplGameInputIllegal
			return
		}
		text.WriteString(fmt.Sprintf(models.TplGameParcelNumber+"\n", game.ParcelNumber(gs.Value)))
	}

	gs.Index = (gs.Index + 1) % len(gs.Players) // switch to the next player.

	text.WriteString(fmt.Sprintf(models.TplGameNext, gs.Players[gs.Index].Name)) // output name of the next player.

	cmd.Output <- text.String()
}

// gameStop clear all fields of GameStore.
func gameStop(cmd *command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	if cmd.Room.GameStore.Game == game.Null {
		cmd.Output <- models.TplGameNull
		return
	}

	stop(cmd.Room.GameStore)

	cmd.Output <- models.TplGameStop
}

// gamePlayers output the list of all participating players.
func gamePlayers(cmd *command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	text := strings.Builder{}
	text.WriteString("玩家列表：")
	for k, v := range cmd.Room.GameStore.Players {
		text.WriteString(fmt.Sprintf("\n%d. %s", k+1, v.Name))
	}

	cmd.Output <- text.String()
}

// gameRank output the game ranking in current room.
func gameRank(cmd *command) {
	rank, err := game.ScoreRank(cmd.Room.ID)
	if err != nil {
		zap.S().Error("get score rank failed: ", err)
		return
	}
	if len(rank) == 0 {
		// check if rank list empty.
		cmd.Output <- models.TplGameRankEmpty
		return
	}

	text := strings.Builder{}
	text.WriteString("游戏排行榜：")
	for k, v := range rank {
		username, err := modules.QueryUsername(v.UID)
		if err != nil {
			zap.S().Error("get game ranking failed: ", err)
			username = "UNKNOWN"
		}
		text.WriteString(fmt.Sprintf("\n%d. %s %d分", k+1, username, v.Score))
	}

	cmd.Output <- text.String()
}

// addScore add score for all the players.
func addScore(roomID int, players []game.Player, score int) {
	for _, p := range players {
		game.UpdateScore(roomID, p.ID, score)
	}
}

func stop(s *game.Store) {
	s.Game = game.Null
	s.State = game.StateNotCreated
	s.Players = nil
	s.Value = nil
}
