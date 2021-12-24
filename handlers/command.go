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

// cmdHandler is the function type that receives *command
// and handle the command event.
type cmdHandler func(cmd *models.Command)

var _cmdMap = map[int]cmdHandler{
	models.CmdBotHelper:   botHelper,
	models.CmdBotFeatures: botUpdates,
	models.CmdRoomInfo:    roomInfo,
	models.CmdCheckin:     checkin,
	models.CmdCheckinRank: checkinRank,
	models.CmdHoroscope:   apiHoroscopes,
	models.CmdWeather:     apiWeather,
	models.CmdSongReq:     songReq,
	models.CmdSongList:    songList,
	models.CmdSongDone:    songDone,
	models.CmdPiaStart:    piaStart,
	models.CmdPiaNext:     piaNext,
	models.CmdPiaNextSafe: piaNextSafe,
	models.CmdPiaRelocate: piaRelocate,
	models.CmdPiaStop:     piaStop,
	models.CmdModeAll:     modeAll,
	models.CmdModeMute:    switchMute,
	models.CmdModePander:  switchPander,
	models.CmdModePinyin:  switchPinyin,
	models.CmdModeWater:   switchWater,
	models.CmdGameRank:    gameRank,
	models.CmdLove:        love,
}

// botHelper outputs the
func botHelper(cmd *models.Command) {
	args := cmd.Args
	key := ""
	if len(args) == 0 {
		key = modules.TmplHelper
	} else {
		switch args[0] {
		case "点歌", "歌单", "req", "list":
			key = modules.TmplHelperPlaylist
		case "pia戏", "pia":
			key = modules.TmplHelperPia
		case "游戏", "game":
			key = modules.TmplHelperGame
		}
	}

	text, err := modules.NewTemplate(key, nil)
	if err != nil {
		zap.S().Warn(cmd.Room.Log("create template failed", err))
		return
	}
	cmd.Output <- text
}

func botUpdates(cmd *models.Command) {
	text, err := modules.NewTemplate(modules.TmplUpdates, nil)
	if err != nil {
		zap.S().Warn(cmd.Room.Log("create template failed", err))
		return
	}
	cmd.Output <- text
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

	text, err := modules.NewTemplate(modules.TmplRoomInfo, data)
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
		zap.S().Warn(cmd.Room.Log("checkin error", err))
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

// The apiHoroscopes handle the apiHoroscopes commands.
func apiHoroscopes(cmd *models.Command) {
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
		zap.S().Warn(cmd.Room.Log("get horoscopes failed", err))
		return
	}
	if n > 0 {
		ret := rdb.HMGet(ctx, key, "content", "score")
		cmd.Output <- fmt.Sprintf(models.TplStarFortune, str, ret.Val()[1], ret.Val()[0])
	} else {
		// check if exists
		fort, err := thirdparty.Zodiac(str, thirdparty.Level5)
		if err != nil {
			zap.S().Warn(cmd.Room.Log("get horoscopes failed", err))
			return
		}
		rdb.HMSet(ctx, key, "content", fort.Content, "score", fort.Score)
		cmd.Output <- fmt.Sprintf(models.TplStarFortune, str, fort.Score, fort.Content)
	}
}

// The apiWeather handles the apiWeather command.
func apiWeather(cmd *models.Command) {
	if len(cmd.Args) != 1 {
		return
	}

	text, err := thirdparty.WeatherText(cmd.Args[0])
	if err != nil {
		zap.S().Warn(cmd.Room.Log("get weather failed", err))
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

func songReq(cmd *models.Command) {
	if len(cmd.Args) != 1 {
		return
	}

	song := cmd.Args[0]
	playlist := cmd.Room.Playlist
	playlist.PushBack(order{song, cmd.User.Username})

	cmd.Output <- fmt.Sprintf(models.TplSongReq, song)
}

func songList(cmd *models.Command) {
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

	text, err := modules.NewTemplate(modules.TmplPlaylist, songs)
	if err != nil {
		zap.S().Warn(cmd.Room.Log("create template failed", err))
		return
	}

	cmd.Output <- text
}

func songDone(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	playlist := cmd.Room.Playlist
	playlist.Remove(playlist.Front())

	cmd.Output <- models.TplSongDone
}

func piaStart(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	if len(cmd.Args) == 0 {
		return
	}

	id, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return
	}

	var roles []string
	roles, cmd.Room.PiaList, err = thirdparty.Fetch(id)
	if err != nil {
		zap.S().Warn(cmd.Room.Log("get scripts failed", err))
		return
	}
	cmd.Room.PiaIndex = 1
	text := strings.Builder{}
	for _, v := range roles {
		text.WriteString(fmt.Sprintf("\n%s", v))
	}
	cmd.Output <- fmt.Sprintf(models.TplPiaStart, text.String())
}

func piaNext(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	args := cmd.Args
	if len(args) == 0 {
		piaNextN(cmd, 1, false)
		return
	}
	dur, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}
	piaNextN(cmd, dur, false)
}

// piaNextSafe outputs one paragraph in safety mode.
func piaNextSafe(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	piaNextN(cmd, 1, true)
}

// piaNextN outputs dur paragraphs.
// If safe=true, output paragraphs in safety mode
// that can avoid being blocked.
func piaNextN(cmd *models.Command, dur int, safe bool) {
	if cmd.Room.PiaIndex == 0 || cmd.Room.PiaList == nil {
		// There is no script in PiaList.
		cmd.Output <- models.TplPiaEmpty
		return
	}

	length := len(cmd.Room.PiaList)
	start := cmd.Room.PiaIndex - 1
	stop := start + dur
	if stop > length {
		// Index out of bounds, relocate to tail of the list.
		stop = length
	}

	text := strings.Builder{}
	for k, v := range (cmd.Room.PiaList)[start:stop] {
		if k > 0 {
			text.WriteString("\n")
		}
		if safe {
			// Add whitespace between each character, avoid being blocked.
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
		// If reaching the end, clear the list and index.
		cmd.Room.PiaList = nil
		cmd.Room.PiaIndex = 0
		cmd.Output <- models.TplPiaDone
	}
}

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

	data := modules.ModeAll(cmd.Room.ID)
	text, err := modules.NewTemplate(modules.TmplMode, data)
	if err != nil {
		zap.S().Warn(cmd.Room.Log("create template failed", err))
		return
	}
	cmd.Output <- text
}

func switchMute(cmd *models.Command) {
	if cmd.Role > models.RoleCreator {
		return
	}
	cmd.Output <- models.TplModeSwitch
	modules.SwitchMode(cmd.Room.ID, modules.ModeMute)
}

func switchPinyin(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	modules.SwitchMode(cmd.Room.ID, modules.ModePinyin)
	cmd.Output <- models.TplModeSwitch
}

func switchPander(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	room := cmd.Room
	args := cmd.Args
	pander := modules.ModePander
	if len(args) >= 1 {
		if !modules.Validate(args[0]) {
			cmd.Output <- models.TplIllegal
			return
		}
		modules.SetMode(room.ID, pander, args[0])
	} else {
		modules.SwitchMode(room.ID, pander)
	}
	cmd.Output <- models.TplModeSwitch
}

func switchWater(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	room := cmd.Room
	args := cmd.Args
	water := modules.ModeWater
	if len(args) >= 1 {
		if !modules.Validate(args[0]) {
			cmd.Output <- models.TplIllegal
			return
		}
		modules.SetMode(room.ID, water, args[0])
	} else {
		modules.SwitchMode(room.ID, water)
	}
	cmd.Output <- models.TplModeSwitch
}

func gameRank(cmd *models.Command) {
	rank, err := modules.ScoreRank(cmd.Room.ID)
	if err != nil {
		zap.S().Warn(cmd.Room.Log("get score rank failed", err))
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
			zap.S().Warn(cmd.Room.Log("get game ranking failed", err))
			username = "UNKNOWN"
		}
		text.WriteString(fmt.Sprintf("\n%d. %s %d分", k+1, username, v.Score))
	}

	cmd.Output <- text.String()
}

func love(cmd *models.Command) {
	cmd.Output <- "❤️~"
}
