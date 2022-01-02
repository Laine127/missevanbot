package models

// command constants.
const (
	CmdBotHelper   = iota // help message
	CmdBotFeatures        // updates of bot recently

	CmdRoomInfo
	CmdRoomAdmin
	CmdCheckin
	CmdCheckinRank
	CmdHoroscope
	CmdWeather

	CmdModeAll
	CmdModeMute // disable message sender of the bot
	CmdModePander
	CmdModePinyin
	CmdModeNoble
	CmdModeMedal
	CmdModeWater // water reminder

	CmdSongReq
	CmdSongList
	CmdSongDone
	CmdSongClear

	CmdPiaStart
	CmdPiaNext
	CmdPiaStop

	CmdGameRank

	CmdLove
)

const (
	RoleSuper   = iota // admin of the bot
	RoleCreator        // room creator
	RoleAdmin          // room admin
	RoleMember         // general member
)

type Command struct {
	*Room
	Args   []string // exclude command
	User   FmUser
	Info   FmInfo
	Role   int
	Output chan<- string
}

var _cmdMap = map[string]int{
	// Bot commands.
	"帮助":   CmdBotHelper,
	"help": CmdBotHelper,
	"更新":   CmdBotFeatures,
	"feat": CmdBotFeatures,
	// Common commands.
	"房间":      CmdRoomInfo,
	"room":    CmdRoomInfo,
	"info":    CmdRoomInfo,
	"管理员":     CmdRoomAdmin,
	"admin":   CmdRoomAdmin,
	"签到":      CmdCheckin,
	"打卡":      CmdCheckin,
	"dd":      CmdCheckin,
	"checkin": CmdCheckin,
	"签到榜":     CmdCheckinRank,
	"星座":      CmdHoroscope,
	"天气":      CmdWeather,
	// Playlist commands.
	"点歌":    CmdSongReq,
	"req":   CmdSongReq,
	"歌单":    CmdSongList,
	"list":  CmdSongList,
	"完成":    CmdSongDone,
	"done":  CmdSongDone,
	"清空":    CmdSongClear,
	"clear": CmdSongClear,
	// Drama commands, no English version.
	"贴本": CmdPiaStart,
	"选本": CmdPiaStart,
	"n":  CmdPiaNext,
	"结束": CmdPiaStop,
	// Game commands.
	"排行":   CmdGameRank,
	"rank": CmdGameRank,
	// Hidden commands.
	"比心":   CmdLove,
	"笔芯":   CmdLove,
	"love": CmdLove,
	// Mode switch commands.
	"模式":     CmdModeAll,
	"mode":   CmdModeAll,
	"静音":     CmdModeMute,
	"mute":   CmdModeMute,
	"显示拼音":   CmdModePinyin,
	"pinyin": CmdModePinyin,
	"显示贵族":   CmdModeNoble,
	"noble":  CmdModeNoble,
	"显示粉丝牌":  CmdModeMedal,
	"medal":  CmdModeMedal,
	"喝水助手":   CmdModeWater,
	"water":  CmdModeWater,
	"咳咳":     CmdModePander,
	"彩虹屁":    CmdModePander,
	"pander": CmdModePander,
}

// Cmd use key to get command constant.
func Cmd(key string) int {
	if v, ok := _cmdMap[key]; ok {
		return v
	}
	return -1
}
