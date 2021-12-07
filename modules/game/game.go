package game

// 游戏种类
const (
	Null       = iota // 无游戏
	NumberBomb        // 数字炸弹
	PassParcel        // 击鼓传花
)

// 游戏状态
const (
	StateNotCreated = iota // 未创建
	StateCreated           // 已创建
	StateReady             // 准备就绪，等待加入
	StateRunning           // 正在进行
)

// 游戏指令
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

// HelpText 帮助文本
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

// _cmdMap 命令映射
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

// Command 使用 Key 获取对应的命令
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
	Game    int            // 游戏种类
	State   int            // 游戏状态
	Players []Player       // 玩家ID列表，按顺序
	Index   int            // 下一个玩家的下标
	Value   map[string]int // 存储数据
}

// AddPlayer 添加玩家
func (s *Store) AddPlayer(player Player) bool {
	for _, v := range s.Players {
		if v.ID == player.ID {
			return false
		}
	}
	s.Players = append(s.Players, player)
	return true
}
