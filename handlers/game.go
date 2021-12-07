package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"missevan-fm/models"
	"missevan-fm/modules"
	"missevan-fm/modules/game"
)

// gameCreate 创建游戏实例，当前状态必须为 game.StateNotCreated 即未创建，
// 需要对存储的 GameStore 进行初始化操作
func gameCreate(cmd *command, gameType int) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	store := cmd.Room.GameStore
	if store.State != game.StateNotCreated {
		// 游戏已经存在了
		cmd.Output <- models.TplGameExists
		return
	}

	switch gameType {
	case game.CmdNumberBomb:
		store.Game = game.NumberBomb
	case game.CmdPassParcel:
		store.Game = game.PassParcel
	default:
		return
	}
	store.State = game.StateCreated // state transfer
	store.Players = make([]game.Player, 0)
	store.Value = make(map[string]int)

	cmd.Output <- models.TplGameCreate
}

// gameJoin 进行玩家加入操作，当前状态必须为 game.StateCreated 即已创建，
// 将玩家 ID 加入玩家列表当中，当数量不小于 2 时切换为 game.StateReady 已就绪状态
func gameJoin(cmd *command, textMsg models.FmTextMessage) {
	store := cmd.Room.GameStore
	if store.State != game.StateCreated && store.State != game.StateReady {
		return
	}

	user := textMsg.User
	if !store.AddPlayer(game.Player{ID: user.UserID, Name: user.Username}) {
		cmd.Output <- models.TplGameJoinDup
		return
	}

	if len(store.Players) >= 2 {
		store.State = game.StateReady
	}

	cmd.Output <- fmt.Sprintf(models.TplGameJoin, user.Username)
}

// gameStart 进行游戏开始操作，当前状态必须为 game.StateReady 即已就绪，
// 进行状态转移变为 game.StateRunning
func gameStart(cmd *command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	store := cmd.Room.GameStore
	if store.State != game.StateReady {
		if len(store.Players) < 2 {
			cmd.Output <- models.TplGameNotEnough
		}
		return
	}

	store.State = game.StateRunning // state transfer
	store.Index = 0                 // 设置第一个玩家

	var min, max int // set for number bomb game
	switch store.Game {
	case game.NumberBomb:
		min, max = game.BombGenerate(store.Value, len(store.Players))
	case game.PassParcel:
		go func() {
			// 定时，每五个玩家设置一分钟
			timer := time.NewTimer(time.Minute * time.Duration(len(store.Players)/5+1))
			<-timer.C
			cmd.Output <- fmt.Sprintf(models.TplGameParcelOver, store.Players[store.Index].Name)
			// 为其他玩家加分
			if store.Index > 0 {
				addScore(cmd.Room.ID, store.Players[:store.Index], 5) // 加分
			}
			if store.Index < len(store.Players)-1 {
				addScore(cmd.Room.ID, store.Players[store.Index+1:], 5) // 加分
			}
			// stop the game
			stop(cmd.Room.GameStore)
			return
		}()
	default:
		return
	}

	cmd.Output <- models.TplGameStart

	switch store.Game {
	case game.NumberBomb:
		cmd.Output <- fmt.Sprintf(models.TplGameBombRange, min, max)
	case game.PassParcel:
		cmd.Output <- fmt.Sprintf(models.TplGameParcelNumber, game.ParcelNumber(store.Value))
	}

	time.Sleep(time.Millisecond * 100)

	cmd.Output <- fmt.Sprintf(models.TplGameNext, store.Players[0].Name)
}

// gameAction 处理游戏中的操作，当前状态必须为 game.StateRunning，
// 操作玩家必须为 Index 指定的玩家
func gameAction(cmd *command, textMsg models.FmTextMessage) {
	store := cmd.Room.GameStore
	if store.State != game.StateRunning {
		return
	}

	if store.Players[store.Index].ID != textMsg.User.UserID {
		// 错误的玩家
		return
	}

	switch store.Game {
	case game.NumberBomb:
		num, err := strconv.Atoi(textMsg.Message)
		if err != nil {
			cmd.Output <- models.TplGameInputIllegal
			return
		}
		ok, min, max := game.BombGuess(store.Value, num)
		if ok {
			cmd.Output <- fmt.Sprintf(models.TplGameBomb, textMsg.User.Username)
			// 为其他玩家加分
			if store.Index > 0 {
				addScore(cmd.Room.ID, store.Players[:store.Index], 10) // 加分
			}
			if store.Index < len(store.Players)-1 {
				addScore(cmd.Room.ID, store.Players[store.Index+1:], 10) // 加分
			}
			// stop the game
			stop(cmd.Room.GameStore)
			return
		}
		if min == -1 {
			cmd.Output <- models.TplGameInputIllegal
			return
		}
		cmd.Output <- fmt.Sprintf(models.TplGameBombRange, min, max)
	case game.PassParcel:
		if !game.IsParcelCorrect(store.Value, textMsg.Message) {
			// 输入不正确
			cmd.Output <- models.TplGameInputIllegal
			return
		}
		cmd.Output <- fmt.Sprintf(models.TplGameParcelNumber, game.ParcelNumber(store.Value))
	default:
		return
	}

	store.Index = (store.Index + 1) % len(store.Players) // change to next player
	cmd.Output <- fmt.Sprintf(models.TplGameNext, store.Players[store.Index].Name)
}

// gameStop 将状态全部归为原值
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

// gamePlayers 输出所有玩家列表
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

// gameRank 输出当前直播间的游戏排行榜
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
			zap.S().Error("get score rank failed: ", err)
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
