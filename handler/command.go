package handler

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"missevan-fm/bot"
	"missevan-fm/module"
)

// 指令类型
const (
	CmdHelper   = iota // 帮助提示
	CmdInfo            // 直播间信息
	CmdSign            // 签到答复
	CmdRank            // 榜单答复
	CmdLove            // 比心答复
	CmdBait            // 演员模式启停
	CmdWeather         // 天气
	CmdMusicAdd        // 点歌
	CmdMusicAll        // 点歌歌单
	CmdMusicPop        // 弹出一首歌
)

// 用户角色
const (
	RoleSuper   = iota // 机器人管理员
	RoleCreator        // 主播
	RoleAdmin          // 房管
	RoleMember         // 普通成员
)

// helpText 帮助文本
const helpText = `命令帮助：
帮助 -- 获取帮助信息
房间 -- 查看当前直播间信息
签到 -- 在当前直播间进行签到
排行 -- 查看当前直播间当天签到排行
天气 城市名 -- 查询该城市的当日天气
点歌 歌名 -- 将一首歌曲添加进待播歌单
歌单 -- 查询当前待播歌单
完成 -- 删除待播清单第一首歌曲`

// _cmdMap 帮助映射
var _cmdMap = map[string]int{
	"帮助": CmdHelper,
	"房间": CmdInfo,
	"签到": CmdSign,
	"排行": CmdRank,
	"天气": CmdWeather,
	"点歌": CmdMusicAdd,
	"歌单": CmdMusicAll,
	"完成": CmdMusicPop,
	// 下面是隐藏的命令
	"比心": CmdLove,
	"笔芯": CmdLove,
	"咳咳": CmdBait,
}

type command struct {
	Room *RoomStore
	User *FmUser
	Role int
}

// info 处理直播间相关的命令
func (cmd *command) info(info *module.Info) {
	if cmd.Role > RoleAdmin {
		return // 权限不足
	}
	text := fmt.Sprintf("当前直播间信息：\n- 房间名：%s\n- 公告：%s\n- 主播：%s\n- 在线人数：%d\n- 累计人数：%d\n- 管理员：\n",
		info.Room.Name,
		info.Room.Announcement,
		info.Creator.Username,
		cmd.Room.Online,
		info.Room.Statistics.Accumulation,
	)
	for k, v := range info.Room.Members.Admin {
		text += fmt.Sprintf("--- %s", v.Username)
		if k < len(info.Room.Members.Admin)-1 {
			text += "\n"
		}
	}
	module.MustSend(info.Room.RoomID, text)
}

// sign 处理签到命令
func (cmd *command) sign(user FmUser) {
	ret, err := module.Sign(cmd.Room.ID, user.UserID, user.Username)
	if err != nil {
		cmd.Room.Error(errors.New("签到出错了：" + err.Error()))
		return
	}
	text := fmt.Sprintf("@%s %s", user.Username, ret)
	module.MustSend(cmd.Room.ID, text)
}

// rank 处理排行榜命令
func (cmd *command) rank() {
	var text string
	if rank := module.Rank(cmd.Room.ID); rank != "" {
		text = fmt.Sprintf("每日签到榜单：%s", rank)
	} else {
		text = "今天的榜单好像空空的~"
	}
	module.MustSend(cmd.Room.ID, text)
}

// bait 处理演员模式启停命令
func (cmd *command) bait() {
	if cmd.Role > RoleAdmin {
		return // 权限不足
	}
	room := cmd.Room
	if room.Bait && room.Timer != nil {
		module.MustSend(room.ID, "我突然有点困了")
		room.Bait = false
		room.Timer.Stop()
	} else {
		room.Bait = true
		if room.RainbowMaxInterval <= 0 {
			room.RainbowMaxInterval = 10
		}
		// 启用定时任务
		timer := time.NewTimer(1)
		room.Timer = timer
		go module.Praise(room.RoomConfig, timer)
	}
}

// weather 处理天气查询命令
func (cmd *command) weather(city string) {
	text := module.Weather(city)
	if text == "" {
		text = "查询的城市好像不正确哦~"
	}
	module.MustSend(cmd.Room.ID, text)
}

// musicAdd 处理点歌命令
func (cmd *command) musicAdd(music string) {
	module.MusicAdd(cmd.Room.ID, music)
	module.MustSend(cmd.Room.ID, fmt.Sprintf("点歌 %s 成功~", music))
}

// musicAll 处理歌单获取命令
func (cmd *command) musicAll() {
	if cmd.Role > RoleAdmin {
		return // 权限不足
	}
	musics := module.MusicAll(cmd.Room.ID)
	if len(musics) == 0 {
		module.MustSend(cmd.Room.ID, "当前还没有人点歌哦~")
		return
	}
	text := strings.Builder{}
	for k, v := range musics {
		text.WriteString(fmt.Sprintf("%d. %s", k+1, v))
		if k < len(musics)-1 {
			text.WriteString("\n")
		}
	}
	module.MustSend(cmd.Room.ID, text.String())
}

// musicPop 处理弹出歌曲命令
func (cmd *command) musicPop() {
	if cmd.Role > RoleAdmin {
		return // 权限不足
	}
	module.MusicPop(cmd.Room.ID)
	module.MustSend(cmd.Room.ID, "完成了一首歌曲~")
}

// userRole 判断当前用户的角色
func userRole(info *module.Info, userID int) int {
	switch userID {
	case bot.Conf.Admin:
		return RoleSuper // 机器人管理员
	case info.Creator.UserID:
		return RoleCreator // 主播
	default:
		// 判断是否是房管
		for _, v := range info.Room.Members.Admin {
			if v.UserID == userID {
				return RoleAdmin
			}
		}
		return RoleMember // 普通用户
	}
}
