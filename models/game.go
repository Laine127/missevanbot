package models

// game state constants.
const (
	StateNotCreated = iota // game is not created.
	StateCreated           // game is created.
	StateReady             // game is ready, waiting to start.
	StateRunning           // game is running.
)

// game command constants.
const (
	CmdGameNumberBomb = iota
	CmdGamePassParcel
	CmdGameGuessWord
	CmdGameJoin
	CmdGameStart
	CmdGameStop
	CmdGamePlayers
)

// game score constants.
const (
	ScoreNumberBomb = 10
	ScorePassParcel = 5
)

// _gameCmdMap is for command mapping.
var _gameCmdMap = map[string]int{
	"数字炸弹": CmdGameNumberBomb,
	"击鼓传花": CmdGamePassParcel,
	"你说我猜": CmdGameGuessWord,
	"加入":   CmdGameJoin,
	"开始":   CmdGameStart,
	"停止":   CmdGameStop,
	"玩家":   CmdGamePlayers,
}

// CmdGame use key to get command constant.
func CmdGame(key string) int {
	if v, ok := _gameCmdMap[key]; ok {
		return v
	}
	return -1
}

type Player struct {
	ID   int
	Name string
}

type Gamer interface {
	// AddPlayer add the player to the list of players,
	// if the player has already joined, return false.
	AddPlayer(player Player) bool
	AllPlayers(cmd *Command)

	// Create initialize a GameStore instance,
	// current state must be StateNotCreated.
	Create(cmd *Command)
	// Join put message sender into the player list,
	// current state must be StateCreated,
	// if the number of players not less than 2 then switch to the state StateReady.
	Join(cmd *Command, textMsg FmTextMessage)
	// Start the game, current state must be StateReady,
	// if there is no error then switch to the state StateRunning.
	Start(cmd *Command)
	// Action handle all game actions, current state must be StateRunning,
	// message sender must be the player pointed to by the index.
	Action(cmd *Command, textMsg FmTextMessage)
	Stop(cmd *Command)
}
