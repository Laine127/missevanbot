package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"missevanbot/config"
	"missevanbot/models"
	"missevanbot/modules"
	"missevanbot/modules/thirdparty"
	"missevanbot/utils"
)

const defParas = 7

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
	models.CmdPiaStop:     piaStop,
	models.CmdModeAll:     modeAll,
	models.CmdModeMute:    switchMute,
	models.CmdModePander:  switchPander,
	models.CmdModePinyin:  switchPinyin,
	models.CmdModeWater:   switchWater,
	models.CmdGameRank:    gameRank,
	models.CmdLove:        love,
}

// botHelper outputs the main help text by default,
// and outputs other help texts according to the arguments.
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
		zap.S().Warn(cmd.Log("create template failed", err))
		return
	}
	cmd.Output <- text
}

func botUpdates(cmd *models.Command) {
	text, err := modules.NewTemplate(modules.TmplUpdates, nil)
	if err != nil {
		zap.S().Warn(cmd.Log("create template failed", err))
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

	openTime := time.UnixMilli(int64(info.Room.Status.OpenTime))
	dur := time.Since(openTime)

	infoRoom := info.Room
	infoCreator := info.Creator

	data := struct {
		Name         string
		Creator      string
		Followers    int
		Platform     string
		Online       int
		Accumulation int
		Count        int
		Admins       []string
		Duration     int
		Vip          int
		Medal        string
	}{
		Name:         infoRoom.Name,
		Creator:      infoCreator.Username,
		Followers:    infoRoom.Statistics.AttentionCount,
		Platform:     infoRoom.Status.Channel.Platform,
		Online:       cmd.Online,
		Accumulation: infoRoom.Statistics.Accumulation,
		Count:        cmd.Count,
		Admins:       admins,
		Duration:     int(dur.Minutes()),
		Vip:          info.Room.Statistics.Vip,
		Medal:        infoRoom.Medal.Name,
	}

	text, err := modules.NewTemplate(modules.TmplRoomInfo, data)
	if err != nil {
		zap.S().Warn(cmd.Log("create template failed", err))
		return
	}
	cmd.Output <- text
}

// The checkin handles the checkin command.
func checkin(cmd *models.Command) {
	user := cmd.User
	uid, name := user.UserID, user.Username
	text, err := modules.Checkin(cmd.ID, uid, name)
	if err != nil {
		zap.S().Warn(cmd.Log("checkin error", err))
		return
	}
	cmd.Output <- text
}

func checkinRank(cmd *models.Command) {
	if rank := modules.CheckinRank(cmd.ID); rank != "" {
		cmd.Output <- fmt.Sprintf("每日签到榜单：%s", rank)
		return
	}
	cmd.Output <- fmt.Sprintf(models.TplRankEmpty, cmd.User.Username)
}

// TODO: find a new API to represent this function.
func apiHoroscopes(cmd *models.Command) {
	args := cmd.Args
	if len(args) == 0 {
		cmd.Output <- fmt.Sprintf(models.TplIllegal, cmd.User.Username)
		return
	}

	str := args[0]
	if len(str) == 6 {
		str += "座"
	}
	if _, ok := thirdparty.StarList[str]; !ok {
		return
	}

	ctx := context.Background()
	rdb := config.RDB
	key := fmt.Sprintf("%szodiac:%s:%s", config.RedisPrefix, utils.Today(), str)

	n, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		zap.S().Warn(cmd.Log("get horoscopes failed", err))
		return
	}
	if n > 0 {
		ret := rdb.HMGet(ctx, key, "content", "score")
		cmd.Output <- fmt.Sprintf(models.TplStarFortune, str, ret.Val()[1], ret.Val()[0])
	} else {
		// check if exists
		fort, err := thirdparty.Zodiac(str, thirdparty.Level5)
		if err != nil {
			zap.S().Warn(cmd.Log("get horoscopes failed", err))
			return
		}
		rdb.HMSet(ctx, key, "content", fort.Content, "score", fort.Score)
		cmd.Output <- fmt.Sprintf(models.TplStarFortune, str, fort.Score, fort.Content)
	}
}

func apiWeather(cmd *models.Command) {
	args := cmd.Args
	if len(args) == 0 {
		cmd.Output <- fmt.Sprintf(models.TplIllegal, cmd.User.Username)
		return
	}

	text, err := thirdparty.WeatherText(args[0])
	if err != nil {
		zap.S().Warn(cmd.Log("get weather failed", err))
		return
	}
	cmd.Output <- text
}

// The order is the struct that used for
// storing the song's name and the command sender's name.
type order struct {
	Name string
	User string
}

func songReq(cmd *models.Command) {
	args := cmd.Args
	if len(args) == 0 {
		cmd.Output <- fmt.Sprintf(models.TplIllegal, cmd.User.Username)
		return
	}

	name := cmd.User.Username
	list := cmd.Playlist
	for _, song := range args {
		list.PushBack(order{song, name})
	}
	cmd.Output <- fmt.Sprintf(models.TplSongReq, name, len(args))
}

func songList(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	list := cmd.Playlist
	ele := list.Front()
	songs := make([]order, 0, 10)
	for i := 0; i < 10 && ele != nil; i++ {
		songs = append(songs, ele.Value.(order))
		ele = ele.Next()
	}

	text, err := modules.NewTemplate(modules.TmplPlaylist, songs)
	if err != nil {
		zap.S().Warn(cmd.Log("create template failed", err))
		return
	}
	cmd.Output <- text
}

func songDone(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	list := cmd.Playlist
	list.Remove(list.Front())
	cmd.Output <- fmt.Sprintf(models.TplSongDone, cmd.User.Username)
}

func piaStart(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}

	args := cmd.Args
	if len(args) == 0 {
		cmd.Output <- fmt.Sprintf(models.TplIllegal, cmd.User.Username)
		return
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}

	var roles []string
	roles, cmd.PiaParas, err = thirdparty.DramaScript(id)
	if err != nil {
		zap.S().Warn(cmd.Log("get scripts failed", err))
		return
	}
	cmd.PiaIndex = 1

	text := strings.Builder{}
	for _, v := range roles {
		text.WriteString(fmt.Sprintf("\n%s", v))
	}
	cmd.Output <- fmt.Sprintf(models.TplPiaStart, cmd.User.Username, text.String())
}

// piaNext outputs N paragraph in safety mode.
func piaNext(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	args := cmd.Args
	if len(args) == 0 {
		piaNextN(cmd, defParas)
		return
	}
	dur, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}
	piaNextN(cmd, dur)
}

// piaNextN outputs dur paragraphs.
// If safe=true, output paragraphs in safety mode
// that can avoid being blocked.
func piaNextN(cmd *models.Command, dur int) {
	index := cmd.PiaIndex
	list := cmd.PiaParas
	if index == 0 || list == nil {
		// There is no script in PiaParas.
		cmd.Output <- fmt.Sprintf(models.TplPiaEmpty, cmd.User.Username)
		return
	}

	length := len(list)
	start := index - 1
	stop := start + dur
	if stop > length {
		// Index out of bounds, relocate to tail of the list.
		stop = length
	}

	text := strings.Builder{}
	for k, v := range (list)[start:stop] {
		if k > 0 {
			text.WriteString("\n")
		}
		// Add zero-width space between each character, avoid being blocked.
		for _, s := range v {
			text.WriteRune(s)
			text.WriteRune('\u200B') // zero-width space
		}
	}
	text.WriteString(fmt.Sprintf("\n\n进度：%d/%d", stop, length))

	cmd.PiaIndex = stop + 1

	if stop == length {
		// If reaching the end, clear the list and index.
		cmd.PiaParas = nil
		cmd.PiaIndex = 0
		text.WriteString("\n\n" + models.TplPiaDone)
	}
	cmd.Output <- text.String()
}

func piaStop(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	if cmd.PiaParas == nil {
		return
	}

	cmd.PiaParas = nil
	cmd.PiaIndex = 0
	cmd.Output <- fmt.Sprintf(models.TplPiaStop, cmd.User.Username)
}

func modeAll(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	data := modules.ModeAll(cmd.ID)
	text, err := modules.NewTemplate(modules.TmplMode, data)
	if err != nil {
		zap.S().Warn(cmd.Log("create template failed", err))
		return
	}
	cmd.Output <- text
}

func switchMute(cmd *models.Command) {
	if cmd.Role > models.RoleCreator {
		return
	}
	cmd.Output <- fmt.Sprintf(models.TplModeSwitch, cmd.User.Username)
	modules.SwitchMode(cmd.ID, modules.ModeMute)
}

func switchPinyin(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	modules.SwitchMode(cmd.ID, modules.ModePinyin)
	cmd.Output <- fmt.Sprintf(models.TplModeSwitch, cmd.User.Username)
}

func switchPander(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	switchValue(cmd, modules.ModePander)
}

func switchWater(cmd *models.Command) {
	if cmd.Role > models.RoleAdmin {
		return
	}
	switchValue(cmd, modules.ModeWater)
}

func switchValue(cmd *models.Command, key string) {
	args := cmd.Args
	name := cmd.User.Username
	if len(args) == 0 {
		modules.SwitchMode(cmd.ID, key)
	} else {
		if !modules.Validate(args[0]) {
			cmd.Output <- fmt.Sprintf(models.TplIllegal, name)
			return
		}
		modules.SetMode(cmd.ID, key, args[0])
	}
	cmd.Output <- fmt.Sprintf(models.TplModeSwitch, name)
}

func gameRank(cmd *models.Command) {
	rank, err := modules.ScoreRank(cmd.ID)
	if err != nil {
		zap.S().Warn(cmd.Log("get score rank failed", err))
		return
	}
	if len(rank) == 0 {
		// Check if rank list empty.
		cmd.Output <- models.TplGameRankEmpty
		return
	}

	text := strings.Builder{}
	text.WriteString("游戏排行榜：")
	for k, v := range rank {
		name, err := modules.QueryUsername(v.UID)
		if err != nil {
			zap.S().Warn(cmd.Log("get game ranking failed", err))
			name = "UNKNOWN"
		}
		text.WriteString(fmt.Sprintf("\n%d. %s %d分", k+1, name, v.Score))
	}
	cmd.Output <- text.String()
}

func love(cmd *models.Command) {
	cmd.Output <- "❤️~"
}
