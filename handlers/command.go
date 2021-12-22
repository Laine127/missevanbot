package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"missevanbot/config"
	"missevanbot/models"
	"missevanbot/modules"
	"missevanbot/modules/thirdparty"
	"missevanbot/utils"
)

// CmdHandler is the function type that receives *command
// and handle the command event.
type CmdHandler func(cmd *models.Command)

var _cmdMap = map[int]CmdHandler{
	models.CmdRoomInfo:    roomInfo,
	models.CmdCheckin:     checkin,
	models.CmdCheckinRank: checkinRank,
	models.CmdHoroscope:   horoscopes,
	models.CmdWeather:     weather,
	models.CmdSongAdd:     songAdd,
	models.CmdSongAll:     songAll,
	models.CmdSongPop:     songPop,
	models.CmdPiaStart:    piaStart,
	models.CmdPiaNext:     piaNext,
	models.CmdPiaNextSafe: piaNextSafe,
	models.CmdPiaRelocate: piaRelocate,
	models.CmdPiaStop:     piaStop,
	models.CmdModeAll:     modeAll,
	models.CmdModeMute:    muteSwitch,
	models.CmdModeBait:    baitSwitch,
	models.CmdModePinyin:  pinyinSwitch,
	models.CmdGameRank:    gameRank,
}

// The roomInfo handles the live room information fetching command.
func roomInfo(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	info := cmd.Info
	fmAdmins := info.Room.Members.Admin

	// The admins slice only contains the name of
	// each live room admin.
	admins := make([]string, 0, len(fmAdmins))
	for _, v := range fmAdmins {
		admins = append(admins, v.Username)
	}

	data := struct {
		Name         string
		Creator      string
		Followers    int
		Platform     string
		Online       int
		Accumulation int
		Count        int
		Admins       []string
	}{
		info.Room.Name,
		info.Creator.Username,
		info.Room.Statistics.AttentionCount,
		info.Room.Status.Channel.Platform,
		cmd.Room.Online,
		info.Room.Statistics.Accumulation,
		cmd.Room.Count,
		admins,
	}

	text, err := models.NewTemplate(models.TmplRoomInfo, data)
	if err != nil {
		zap.Error(err)
		return
	}

	cmd.Output <- text
}

// The checkin handles the checkin command.
func checkin(cmd *models.Command) {
	user := cmd.User
	ret, err := modules.Checkin(cmd.Room.ID, user.UserID, user.Username)
	if err != nil {
		zap.S().Error("checkin error: ", err)
		return
	}

	cmd.Output <- fmt.Sprintf("@%s %s", user.Username, ret)
}

// The checkinRank handles the checkin-rank querying command.
func checkinRank(cmd *models.Command) {
	var text string
	if rank := modules.CheckinRank(cmd.Room.ID); rank != "" {
		text = fmt.Sprintf("每日签到榜单：%s", rank)
	} else {
		text = models.TplRankEmpty
	}

	cmd.Output <- text
}

// The horoscopes handle the horoscopes commands.
func horoscopes(cmd *models.Command) {
	if len(cmd.Args) != 1 {
		return
	}

	str := cmd.Args[0]
	if len(str) == 6 {
		str += "座"
	}

	if _, ok := thirdparty.StarList[str]; !ok {
		// check if validate
		return
	}

	ctx := context.Background()
	rdb := config.RDB
	key := fmt.Sprintf("%szodiac:%s:%s", config.RedisPrefix, utils.Today(), str)

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

// The weather handles the weather command.
func weather(cmd *models.Command) {
	if len(cmd.Args) != 1 {
		return
	}

	text, err := thirdparty.WeatherText(cmd.Args[0])
	if err != nil {
		zap.S().Error("天气查询失败：", err)
		return
	}

	cmd.Output <- text
}

// The song is the struct that used for
// storing the song's name and the command sender's name.
type order struct {
	Name string
	User string
}

func songAdd(cmd *models.Command) {
	if len(cmd.Args) != 1 {
		return
	}

	song := cmd.Args[0]
	playlist := cmd.Room.Playlist
	playlist.PushBack(order{song, cmd.User.Username})

	cmd.Output <- fmt.Sprintf(models.TplSongAdd, song)
}

func songAll(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	playlist := cmd.Room.Playlist
	s := playlist.Front()
	songs := make([]order, 0, 10)
	for i := 0; i < 10 && s != nil; i++ {
		songs = append(songs, s.Value.(order))
		s = s.Next()
	}

	text, err := models.NewTemplate(models.TmplPlaylist, songs)
	if err != nil {
		zap.S().Error(err)
		return
	}

	cmd.Output <- text
}

func songPop(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	playlist := cmd.Room.Playlist
	playlist.Remove(playlist.Front()) // pop the first song in playlist.

	cmd.Output <- models.TplSongDone
}

func piaStart(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	// check args
	if len(cmd.Args) != 1 {
		return
	}

	id, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return
	}

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
func piaNext(cmd *models.Command) {
	switch len(cmd.Args) {
	case 0:
		piaNextN(cmd, 1, false)
	case 1:
		dur, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return
		}
		piaNextN(cmd, dur, false)
	}
}

// piaNextSafe 安全输出一条戏文
func piaNextSafe(cmd *models.Command) {
	piaNextN(cmd, 1, true)
}

// piaNextN 进行发送多条戏本文本的处理
func piaNextN(cmd *models.Command, dur int, safe bool) {
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
		stop = length // 索引越界，则边界定位到数组末尾
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

// piaRelocate 重定位文章位置
func piaRelocate(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}

	if len(cmd.Args) != 1 {
		return
	}

	idx, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return
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
func piaStop(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return // 权限不足
	}
	if cmd.Room.PiaList == nil {
		return
	}
	cmd.Room.PiaList = nil
	cmd.Room.PiaIndex = 0
	cmd.Output <- models.TplPiaStop
}

func modeAll(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	text, err := models.NewTemplate(models.TmplModes, modules.ModeAll(cmd.Room.ID))
	if err != nil {
		zap.S().Error(err)
		return
	}
	cmd.Output <- text
}

func muteSwitch(cmd *models.Command) {
	if cmd.Role > models.RoleCreator {
		return
	}

	modules.SwitchMode(cmd.Room.ID, modules.ModeMute)

	cmd.Output <- models.TplModeSwitch
}

// baitSwitch 处理演员模式启停命令
func baitSwitch(cmd *models.Command) {
	if cmd.Role > models.RoleCreator {
		return
	}

	modules.SwitchMode(cmd.Room.ID, modules.ModeBait)

	cmd.Output <- models.TplModeSwitch
}

func pinyinSwitch(cmd *models.Command) {
	if cmd.Role > models.RoleCreator {
		return
	}

	modules.SwitchMode(cmd.Room.ID, modules.ModePinyin)

	cmd.Output <- models.TplModeSwitch
}

func gameRank(cmd *models.Command) {
	rank, err := modules.ScoreRank(cmd.Room.ID)
	if err != nil {
		zap.S().Error("get score rank failed: ", err)
		return
	}
	if len(rank) == 0 {
		// check if rank list empty.
		cmd.Output <- models.TplGameRankEmpty
		return
	}

	text := strings.Builder{}
	text.WriteString("游戏排行榜：")
	for k, v := range rank {
		username, err := modules.QueryUsername(v.UID)
		if err != nil {
			zap.S().Error("get game ranking failed: ", err)
			username = "UNKNOWN"
		}
		text.WriteString(fmt.Sprintf("\n%d. %s %d分", k+1, username, v.Score))
	}

	cmd.Output <- text.String()
}
