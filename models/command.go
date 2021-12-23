package models

// command constants.
const (
	CmdHelper      = iota // 帮助提示
	CmdRoomInfo           // 直播间信息
	CmdCheckin            // 签到答复
	CmdCheckinRank        // 榜单答复
	CmdHoroscope          // 星座运势
	CmdWeather            // 天气
	CmdModeAll            // 全部模式
	CmdModeMute           // 静音模式
	CmdModeBait           // 演员模式启停
	CmdModePinyin         // 拼音模式启停
	CmdModeWater          // 喝水助手启停
	CmdSongAdd            // 点歌
	CmdSongAll            // 点歌歌单
	CmdSongPop            // 弹出一首歌
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
	Args   []string // exclude command.
	Room   *Room
	User   FmUser
	Info   FmInfo
	Role   int
	Output chan<- string
}

var _cmdMap = map[string]int{
	"帮助":      CmdHelper,
	"help":    CmdHelper,
	"房间":      CmdRoomInfo,
	"room":    CmdRoomInfo,
	"签到":      CmdCheckin,
	"打卡":      CmdCheckin,
	"checkin": CmdCheckin,
	"签到榜":     CmdCheckinRank,
	"星座":      CmdHoroscope,
	"天气":      CmdWeather,
	"点歌":      CmdSongAdd,
	"add":     CmdSongAdd,
	"歌单":      CmdSongAll,
	"完成":      CmdSongPop,
	"pop":     CmdSongPop,
	"贴本":      CmdPiaStart,
	"n":       CmdPiaNext,
	"s":       CmdPiaNextSafe,
	"r":       CmdPiaRelocate,
	"结束":      CmdPiaStop,
	"排行":      CmdGameRank,
	"rank":    CmdGameRank,
	// hidden commands below.
	"比心":   CmdLove,
	"笔芯":   CmdLove,
	"love": CmdLove,
	// mode switch
	"mode": CmdModeAll,
	"mute": CmdModeMute,
	"拼音":   CmdModePinyin,
	"多喝热水": CmdModeWater,
	"咳咳":   CmdModeBait,
}

// Cmd use key to get command constant.
func Cmd(key string) int {
	if v, ok := _cmdMap[key]; ok {
		return v
	}
	return -1
}
