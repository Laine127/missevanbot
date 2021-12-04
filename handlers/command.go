package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"go.uber.org/zap"
	"missevan-fm/config"
	"missevan-fm/models"
	"missevan-fm/modules"
	"missevan-fm/modules/thirdparty"
	"missevan-fm/utils"
)

type command struct {
	Room   *models.Room
	User   *models.FmUser
	Role   int
	Output chan<- string
}

// roomInfo 处理直播间相关的命令
func (cmd *command) roomInfo(info *models.FmInfo) {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}
	text := fmt.Sprintf(models.TplRoomInfo,
		info.Room.Name,
		info.Creator.Username,
		info.Room.Status.Channel.Platform,
		cmd.Room.Online,
		info.Room.Statistics.Accumulation,
	)
	for _, v := range info.Room.Members.Admin {
		text += fmt.Sprintf("\n--- %s", v.Username)
	}
	cmd.Output <- text
}

// checkin 处理签到命令
func (cmd *command) checkin(user models.FmUser) {
	ret, err := modules.Checkin(cmd.Room.ID, user.UserID, user.Username)
	if err != nil {
		zap.S().Errorf("签到出错了：%s", err)
		return
	}
	cmd.Output <- fmt.Sprintf("@%s %s", user.Username, ret)
}

// checkinRank 处理排行榜命令
func (cmd *command) checkinRank() {
	var text string
	if rank := modules.CheckinRank(cmd.Room.ID); rank != "" {
		text = fmt.Sprintf("每日签到榜单：%s", rank)
	} else {
		text = models.TplRankEmpty
	}
	cmd.Output <- text
}

// horoscopes 生成星座运势
func (cmd *command) horoscopes(str string) {
	if utf8.RuneCountInString(str) == 2 {
		str += "座"
	}
	if _, ok := thirdparty.StarList[str]; !ok {
		// check if validate
		return
	}
	ctx := context.Background()
	rdb := config.RDB
	key := fmt.Sprintf("%szodiac:%s:%s", models.RedisPrefix, utils.Today(), str)

	n, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		zap.S().Error("获取星座运势错误：", err)
		return
	}
	if n > 0 {
		ret := rdb.HMGet(ctx, key, "content", "score")
		cmd.Output <- fmt.Sprintf(models.TplStarFortune, str, ret.Val()[1], ret.Val()[0])
	} else {
		// check if exists
		fort, err := thirdparty.Zodiac(str, thirdparty.Level5)
		if err != nil {
			zap.S().Error("获取星座运势错误：", err)
			return
		}
		rdb.HMSet(ctx, key, "content", fort.Content, "score", fort.Score)
		cmd.Output <- fmt.Sprintf(models.TplStarFortune, str, fort.Score, fort.Content)
	}
}

// baitSwitch 处理演员模式启停命令
func (cmd *command) baitSwitch() {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}
	room := cmd.Room
	if room.Bait && room.Timer != nil {
		// 演员模式已经启动了，执行关闭
		cmd.Output <- models.TplBaitStop
		room.Bait = false
		room.Timer.Stop()
	} else {
		// 演员模式还未启动，执行启动
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
	text, err := thirdparty.WeatherText(city)
	if err != nil {
		zap.S().Error("天气查询失败：", err)
		return
	}
	cmd.Output <- text
}

// musicAdd 处理点歌命令
func (cmd *command) musicAdd(music string) {
	modules.MusicAdd(cmd.Room.ID, music)
	cmd.Output <- fmt.Sprintf(models.TplMusicAdd, music)
}

// musicAll 处理歌单获取命令
func (cmd *command) musicAll() {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}
	musics := modules.MusicAll(cmd.Room.ID)
	if len(musics) == 0 {
		cmd.Output <- models.TplMusicNone
		return
	}
	text := strings.Builder{}
	for k, v := range musics {
		if k > 0 {
			text.WriteString("\n")
		}
		text.WriteString(fmt.Sprintf("%d. %s", k+1, v))
	}
	cmd.Output <- text.String()
}

// musicPop 处理弹出歌曲命令
func (cmd *command) musicPop() {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}
	modules.MusicPop(cmd.Room.ID)
	cmd.Output <- models.TplMusicDone
}

// piaStart 处理开启pia戏命令
func (cmd *command) piaStart(id int) {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}
	var err error
	var roles []string
	roles, cmd.Room.PiaList, err = thirdparty.Fetch(id)
	if err != nil {
		zap.S().Error("获取戏文出错了：", err)
		return
	}
	cmd.Room.PiaIndex = 1
	text := strings.Builder{}
	for _, v := range roles {
		text.WriteString(fmt.Sprintf("\n%s", v))
	}
	cmd.Output <- fmt.Sprintf(models.TplPiaStart, text.String())
}

// piaNext 处理下一条戏文命令
func (cmd *command) piaNext(dur int, safe bool) {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}
	if cmd.Room.PiaIndex == 0 || cmd.Room.PiaList == nil {
		// 没有启动
		cmd.Output <- models.TplPiaEmpty
		return
	}

	length := len(cmd.Room.PiaList)
	start := cmd.Room.PiaIndex - 1
	stop := start + dur

	if stop > length {
		// 索引越界，则边界定位到数组末尾
		stop = length
	}

	text := strings.Builder{}
	for k, v := range (cmd.Room.PiaList)[start:stop] {
		if k > 0 {
			text.WriteString("\n")
		}
		if safe {
			// 安全输出，防止屏蔽
			for _, s := range v {
				text.WriteString(string(s))
				text.WriteString(" ")
			}
			continue
		}
		text.WriteString(v)
	}
	text.WriteString(fmt.Sprintf("\n\n进度：%d/%d", stop, length))
	cmd.Room.PiaIndex = stop + 1
	cmd.Output <- text.String()

	if stop == length {
		// 到达末尾，清空列表，关闭当前模式
		cmd.Room.PiaList = nil
		cmd.Room.PiaIndex = 0
		cmd.Output <- models.TplPiaDone
	}
}

// piaNextSafe 安全输出一条戏文
func (cmd *command) piaNextSafe() {
	cmd.piaNext(1, true)
}

// piaRelocate 重定位文章位置
func (cmd *command) piaRelocate(idx int) {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}
	if cmd.Room.PiaIndex == 0 || cmd.Room.PiaList == nil {
		// 没有启动
		cmd.Output <- models.TplPiaEmpty
		return
	}
	idx-- // 转换为程序下标

	length := len(cmd.Room.PiaList)
	if idx >= length {
		// 越界错误
		cmd.Output <- models.TplPiaOutBounds
		return
	}
	cmd.Room.PiaIndex = idx + 1
	cmd.Output <- models.TplPiaRelocate
}

// piaStop 处理停止pia戏命令
func (cmd *command) piaStop() {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}
	cmd.Room.PiaList = nil
	cmd.Room.PiaIndex = 0
	cmd.Output <- models.TplPiaStop
}

// userRole 判断当前用户的角色
func userRole(info *models.FmInfo, userID int) int {
	switch userID {
	case config.Admin():
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
