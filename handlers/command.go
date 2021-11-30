package handlers

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"missevan-fm/config"
	"missevan-fm/models"
	"missevan-fm/modules"
	"missevan-fm/modules/thirdparty"
)

type command struct {
	Room   *models.Room
	User   *models.FmUser
	Role   int
	Output chan<- string
}

// info 处理直播间相关的命令
func (cmd *command) info(info *modules.Info) {
	if cmd.Role > models.RoleAdmin {
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
	cmd.Output <- text
}

// sign 处理签到命令
func (cmd *command) sign(user models.FmUser) {
	ret, err := modules.Sign(cmd.Room.ID, user.UserID, user.Username)
	if err != nil {
		zap.S().Errorf("签到出错了：%s", err)
		return
	}
	cmd.Output <- fmt.Sprintf("@%s %s", user.Username, ret)
}

// rank 处理排行榜命令
func (cmd *command) rank() {
	var text string
	if rank := modules.Rank(cmd.Room.ID); rank != "" {
		text = fmt.Sprintf("每日签到榜单：%s", rank)
	} else {
		text = "今天的榜单好像空空的~"
	}
	cmd.Output <- text
}

// bait 处理演员模式启停命令
func (cmd *command) bait() {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}
	room := cmd.Room
	if room.Bait && room.Timer != nil {
		cmd.Output <- "我突然有点困了"
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
		go modules.Praise(cmd.Output, room.RoomConfig, timer)
	}
}

// weather 处理天气查询命令
func (cmd *command) weather(city string) {
	text := thirdparty.WeatherText(city)
	if text == "" {
		text = "查询的城市好像不正确哦~"
	}
	cmd.Output <- text
}

// musicAdd 处理点歌命令
func (cmd *command) musicAdd(music string) {
	modules.MusicAdd(cmd.Room.ID, music)
	cmd.Output <- fmt.Sprintf("点歌 %s 成功~", music)
}

// musicAll 处理歌单获取命令
func (cmd *command) musicAll() {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}
	musics := modules.MusicAll(cmd.Room.ID)
	if len(musics) == 0 {
		cmd.Output <- "当前还没有人点歌哦~"
		return
	}
	text := strings.Builder{}
	for k, v := range musics {
		text.WriteString(fmt.Sprintf("%d. %s", k+1, v))
		if k < len(musics)-1 {
			text.WriteString("\n")
		}
	}
	cmd.Output <- text.String()
}

// musicPop 处理弹出歌曲命令
func (cmd *command) musicPop() {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}
	modules.MusicPop(cmd.Room.ID)
	cmd.Output <- "完成了一首歌曲~"
}

// userRole 判断当前用户的角色
func userRole(info *modules.Info, userID int) int {
	switch userID {
	case config.Config().Admin:
		return models.RoleSuper // 机器人管理员
	case info.Creator.UserID:
		return models.RoleCreator // 主播
	default:
		// 判断是否是房管
		for _, v := range info.Room.Members.Admin {
			if v.UserID == userID {
				return models.RoleAdmin
			}
		}
		return models.RoleMember // 普通用户
	}
}
