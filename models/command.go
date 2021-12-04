package models

// 指令类型
const (
	CmdHelper   = iota // 帮助提示
	CmdInfo            // 直播间信息
	CmdSign            // 签到答复
	CmdRank            // 榜单答复
	CmdStar            // 星座运势
	CmdLove            // 比心答复
	CmdBait            // 演员模式启停
	CmdWeather         // 天气
	CmdMusicAdd        // 点歌
	CmdMusicAll        // 点歌歌单
	CmdMusicPop        // 弹出一首歌
	CmdPiaStart        // 启动pia戏模式
	CmdPiaNext         // 下一条
	CmdPiaStop         // 结束pia戏模式
)

// 用户角色
const (
	RoleSuper   = iota // 机器人管理员
	RoleCreator        // 主播
	RoleAdmin          // 房管
	RoleMember         // 普通成员
)

// HelpText 帮助文本
const HelpText = `命令帮助：
帮助 -- 获取帮助信息
房间 -- 查看当前直播间信息
签到 -- 在当前直播间进行签到
排行 -- 查看当前直播间当天签到排行
星座 星座名 -- 查看星座今日运势
天气 城市名 -- 查询该城市的当日天气
点歌 歌名 -- 将歌曲添加进待播歌单
歌单 -- 查询当前待播歌单
完成 -- 删除待播清单第一首歌曲
本子 ID -- 获取戏文，开启pia戏模式
n -- 下一条内容
n 数字 -- 多条内容
结束 -- 结束pia戏模式

Author: Secriy`

// _cmdMap 帮助映射
var _cmdMap = map[string]int{
	"帮助": CmdHelper,
	"房间": CmdInfo,
	"签到": CmdSign,
	"排行": CmdRank,
	"星座": CmdStar,
	"天气": CmdWeather,
	"点歌": CmdMusicAdd,
	"歌单": CmdMusicAll,
	"完成": CmdMusicPop,
	"本子": CmdPiaStart,
	"n":  CmdPiaNext,
	"结束": CmdPiaStop,
	// 下面是隐藏的命令
	"比心": CmdLove,
	"笔芯": CmdLove,
	"咳咳": CmdBait,
}

// Command 使用 Key 获取对应的命令
func Command(key string) int {
	if v, ok := _cmdMap[key]; ok {
		return v
	}
	return -1
}
