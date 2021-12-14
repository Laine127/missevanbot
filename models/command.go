package models

// command constants.
const (
	CmdHelper      = iota // 帮助提示
	CmdRoomInfo           // 直播间信息
	CmdCheckin            // 签到答复
	CmdCheckinRank        // 榜单答复
	CmdHoroscope          // 星座运势
	CmdBaitSwitch         // 演员模式启停
	CmdWeather            // 天气
	CmdMusicAdd           // 点歌
	CmdMusicAll           // 点歌歌单
	CmdMusicPop           // 弹出一首歌
	CmdPiaStart           // 启动pia戏模式
	CmdPiaNext            // 下一条
	CmdPiaNextSafe        // 下一条（安全输出）
	CmdPiaRelocate        // 重定位
	CmdPiaStop            // 结束pia戏模式
	CmdGameRank           // 游戏排行
	CmdLove               // 比心答复
)

const (
	RoleSuper   = iota // admin of the bot
	RoleCreator        // room creator
	RoleAdmin          // room admin
	RoleMember         // general member
)

type Command struct {
	Args   []string
	Room   *Room
	User   FmUser
	Info   FmInfo
	Role   int
	Output chan<- string
}

// _cmdMap is for command mapping.
var _cmdMap = map[string]int{
	"帮助":  CmdHelper,
	"房间":  CmdRoomInfo,
	"签到":  CmdCheckin,
	"签到榜": CmdCheckinRank,
	"星座":  CmdHoroscope,
	"天气":  CmdWeather,
	"点歌":  CmdMusicAdd,
	"歌单":  CmdMusicAll,
	"完成":  CmdMusicPop,
	"贴本":  CmdPiaStart,
	"n":   CmdPiaNext,
	"s":   CmdPiaNextSafe,
	"r":   CmdPiaRelocate,
	"结束":  CmdPiaStop,
	"排行":  CmdGameRank,
	// hidden commands
	"比心": CmdLove,
	"笔芯": CmdLove,
	"咳咳": CmdBaitSwitch,
}

// Cmd use key to get command constant.
func Cmd(key string) int {
	if v, ok := _cmdMap[key]; ok {
		return v
	}
	return -1
}
