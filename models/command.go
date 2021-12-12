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
	CmdLove               // 比心答复
)

const (
	RoleSuper   = iota // admin of the bot
	RoleCreator        // room creator
	RoleAdmin          // room admin
	RoleMember         // general member
)

const HelpText = `命令帮助：

帮助 -- 获取帮助信息
游戏 -- 获取游戏帮助信息
房间 -- 查询直播间信息
签到 -- 在当前直播间签到
排行 -- 查看当天签到排行
星座 星座名 -- 查询星座今日运势
天气 城市名 -- 查询该城市的当日天气

+++++ 点歌相关 +++++
点歌 歌名 -- 添加歌曲到排队歌单
歌单 -- 查询排队歌单
完成 -- 删除排队清单第一首歌曲

+++++ pia戏相关 +++++
本子 本号 -- 获取戏文，开启pia戏模式
s -- 下一条文本（防屏蔽版）
n -- 下一条文本
n 数字 -- 多条文本
r 数字 -- 定位到指定的位置
结束 -- 结束pia戏模式

Author: Secriy`

// _cmdMap is for command mapping.
var _cmdMap = map[string]int{
	"帮助": CmdHelper,
	"房间": CmdRoomInfo,
	"签到": CmdCheckin,
	"排行": CmdCheckinRank,
	"星座": CmdHoroscope,
	"天气": CmdWeather,
	"点歌": CmdMusicAdd,
	"歌单": CmdMusicAll,
	"完成": CmdMusicPop,
	"本子": CmdPiaStart,
	"n":  CmdPiaNext,
	"s":  CmdPiaNextSafe,
	"r":  CmdPiaRelocate,
	"结束": CmdPiaStop,
	// hidden commands
	"比心": CmdLove,
	"笔芯": CmdLove,
	"咳咳": CmdBaitSwitch,
}

// Command use key to get command constant.
func Command(key string) int {
	if v, ok := _cmdMap[key]; ok {
		return v
	}
	return -1
}
