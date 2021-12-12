package game

// game type constants.
const (
	Null       = iota // no game is set.
	NumberBomb        // game Number Bomb.
	PassParcel        // game Pass the Parcel.
)

// game state constants.
const (
	StateNotCreated = iota // game is not created.
	StateCreated           // game is created.
	StateReady             // game is ready, waiting to start.
	StateRunning           // game is running.
)

// game command constants.
const (
	CmdHelper = iota
	CmdNumberBomb
	CmdPassParcel
	CmdJoin
	CmdStart
	CmdStop
	CmdPlayers
	CmdRank
)

// game score constants.
const (
	ScoreNumberBomb = 10
	ScorePassParcel = 5
)

const HelpText = `游戏帮助：

游戏 -- 获取游戏相关帮助信息

数字炸弹 -- 创建数字炸弹游戏房间
击鼓传花 -- 创建击鼓传花游戏房间

加入游戏 -- 加入当前创建的游戏
开始游戏 -- 开始当前选定的游戏
结束游戏 -- 结束当前进行的游戏
玩家 -- 查看当前房间中的玩家
战绩 -- 查看当前房间游戏战绩排行

Author: Secriy`

// _cmdMap is for command mapping.
var _cmdMap = map[string]int{
	"游戏":   CmdHelper,
	"数字炸弹": CmdNumberBomb,
	"击鼓传花": CmdPassParcel,
	"加入游戏": CmdJoin,
	"开始游戏": CmdStart,
	"结束游戏": CmdStop,
	"玩家":   CmdPlayers,
	"战绩":   CmdRank,
}

// Command use key to get command constant.
func Command(key string) int {
	if v, ok := _cmdMap[key]; ok {
		return v
	}
	return -1
}

type Player struct {
	ID   int
	Name string
}

type Store struct {
	Game    int            // game type.
	State   int            // game state.
	Players []Player       // store a list of players in order, contains ID and name of each player.
	Index   int            // the index of next player.
	Value   map[string]int // used to store some game values of type int.
}

// AddPlayer add the player to the list of players,
// if the player has already joined, return false.
func (s *Store) AddPlayer(player Player) bool {
	// check if the player has joined the game
	for _, v := range s.Players {
		if v.ID == player.ID {
			return false
		}
	}

	s.Players = append(s.Players, player)
	return true
}
