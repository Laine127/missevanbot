package models

// command constants.
const (
	CmdBotHelper   = iota // help message
	CmdBotFeatures        // updates of bot recently

	CmdRoomInfo
	CmdCheckin
	CmdCheckinRank
	CmdHoroscope
	CmdWeather

	CmdModeAll
	CmdModeMute // disable message sender of the bot
	CmdModePander
	CmdModePinyin
	CmdModeWater // water reminder

	CmdSongReq
	CmdSongList
	CmdSongDone

	CmdPiaStart
	CmdPiaNext
	CmdPiaNextSafe
	CmdPiaRelocate
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
	Args   []string // exclude command
	Room   *Room
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
	"签到":      CmdCheckin,
	"打卡":      CmdCheckin,
	"checkin": CmdCheckin,
	"签到榜":     CmdCheckinRank,
	"星座":      CmdHoroscope,
	"天气":      CmdWeather,
	// Playlist commands.
	"点歌":   CmdSongReq,
	"req":  CmdSongReq,
	"歌单":   CmdSongList,
	"list": CmdSongList,
	"完成":   CmdSongDone,
	"done": CmdSongDone,
	// Drama commands, no English version.
	"贴本": CmdPiaStart,
	"选本": CmdPiaStart,
	"n":  CmdPiaNext,
	"s":  CmdPiaNextSafe,
	"r":  CmdPiaRelocate,
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
	"拼音":     CmdModePinyin,
	"pinyin": CmdModePinyin,
	"多喝热水":   CmdModeWater,
	"water":  CmdModeWater,
	"咳咳":     CmdModePander,
	"pander": CmdModePander,
}

// Cmd use key to get command constant.
func Cmd(key string) int {
	if v, ok := _cmdMap[key]; ok {
		return v
	}
	return -1
}
