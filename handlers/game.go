package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"missevan-fm/models"
	"missevan-fm/modules/game"
)

// gameCreate 创建游戏实例，当前状态必须为 game.StateNotCreated 即未创建，
// 需要对存储的 GameStore 进行初始化操作
func gameCreate(cmd *command, gameType int, textMsg models.FmTextMessage) {
	if cmd.Role > models.RoleCreator {
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
	default:
		return
	}
	store.State = game.StateCreated // state transfer
	store.Players = make([]string, 0)
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

	store.AddPlayer(textMsg.User.Username)

	if len(store.Players) >= 2 {
		store.State = game.StateReady
	}

	cmd.Output <- fmt.Sprintf(models.TplGameJoin, textMsg.User.Username)
}

// gameStart 进行游戏开始操作，当前状态必须为 game.StateReady 即已就绪，
// 进行状态转移变为 game.StateRunning
func gameStart(cmd *command, textMsg models.FmTextMessage) {
	if cmd.Role > models.RoleCreator {
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
	store.Player = 0                // 设置第一个玩家

	switch store.Game {
	case game.NumberBomb:
		game.BombGenerate(store.Value, len(store.Players))
	default:
		return
	}

	cmd.Output <- models.TplGameStart
	cmd.Output <- fmt.Sprintf(models.TplGameNext, store.Players[0])
}

// gameAction 处理游戏中的操作，当前状态必须为 game.StateRunning，
// 操作玩家必须为 Player 指定的玩家
func gameAction(cmd *command, textMsg models.FmTextMessage) {
	store := cmd.Room.GameStore
	if store.State != game.StateRunning {
		return
	}

	if store.Players[store.Player] != textMsg.User.Username {
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
			// stop the game
			s := cmd.Room.GameStore
			s.Game = game.Null
			s.State = game.StateNotCreated
			s.Players = nil
			s.Value = nil
			return
		}
		if min == -1 {
			cmd.Output <- models.TplGameInputIllegal
			return
		}
		cmd.Output <- fmt.Sprintf(models.TplGameBombRange, min, max)
	default:
		return
	}

	store.Player = (store.Player + 1) % len(store.Players) // change to next player
	cmd.Output <- fmt.Sprintf(models.TplGameNext, store.Players[store.Player])
}

// gameStop 将状态全部归为原值
func gameStop(cmd *command, textMsg models.FmTextMessage) {
	if cmd.Role > models.RoleCreator {
		return
	}

	store := cmd.Room.GameStore
	store.Game = game.Null
	store.State = game.StateNotCreated
	store.Players = nil
	store.Value = nil

	cmd.Output <- models.TplGameStop
}

// gamePlayers 输出所有玩家列表
func gamePlayers(cmd *command, textMsg models.FmTextMessage) {
	if cmd.Role > models.RoleCreator {
		return
	}

	text := strings.Builder{}
	text.WriteString("玩家列表：")
	for k, v := range cmd.Room.GameStore.Players {
		text.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
	}

	cmd.Output <- text.String()
}
