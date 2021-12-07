package game

const (
	Null       = iota // 无游戏
	NumberBomb        // 数字炸弹
)

const (
	StateNotCreated = iota // 未创建
	StateCreated           // 已创建
	StateReady             // 准备就绪，等待加入
	StateRunning           // 正在进行
)

const (
	CmdHelper = iota
	CmdNumberBomb
	CmdJoin
	CmdStart
	CmdStop
	CmdPlayers
)

// HelpText 帮助文本
const HelpText = `游戏帮助：

游戏 -- 获取游戏相关帮助信息
数字炸弹 -- 创建数字炸弹游戏房间
加入游戏 -- 加入当前创建的游戏
开始游戏 -- 开始当前选定的游戏
终止游戏 -- 终止当前进行的游戏
查看玩家 -- 查看当前房间中的玩家

Author: Secriy`

// _cmdMap 命令映射
var _cmdMap = map[string]int{
	"游戏":   CmdHelper,
	"数字炸弹": CmdNumberBomb,
	"加入游戏": CmdJoin,
	"开始游戏": CmdStart,
	"终止游戏": CmdStop,
	"查看玩家": CmdPlayers,
	// "战绩":   CmdGameResult,
}

// Command 使用 Key 获取对应的命令
func Command(key string) int {
	if v, ok := _cmdMap[key]; ok {
		return v
	}
	return -1
}

type Store struct {
	Game    int            // 游戏种类
	State   int            // 游戏状态
	Players []string       // 玩家ID列表，按顺序
	Player  int            // 下一个玩家的下标
	Value   map[string]int // 存储数据
}

// AddPlayer 添加玩家
func (s *Store) AddPlayer(player string) {
	s.Players = append(s.Players, player)
}
