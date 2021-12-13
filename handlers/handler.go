package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"missevan-fm/config"
	"missevan-fm/models"
	"missevan-fm/modules"
	"missevan-fm/utils"
)

// HandleRoom 处理直播间相关事件
func HandleRoom(outputMsg chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	info, err := modules.RoomInfo(room.ID)
	if err != nil {
		zap.S().Error("获取直播间信息错误：", err)
		return
	}

	switch textMsg.Event {
	case models.EventStatistic:
		// 更新在线人数
		room.Online = textMsg.Statistics.Online
	case models.EventOpen:
		outputMsg <- models.TplBotStart
		// 通知推送
		if room.Watch {
			text := fmt.Sprintf("%s 开播啦~", info.Creator.Username)
			if err := modules.Push(modules.TitleOpen, text); err != nil {
				zap.S().Error("Bark 推送失败", err)
			}
		}
	case models.EventClose:
		// 通知推送
		if room.Watch {
			text := fmt.Sprintf("%s 下播啦~", info.Creator.Username)
			if err := modules.Push(modules.TitleClose, text); err != nil {
				zap.S().Error("Bark 推送失败", err)
			}
		}
		// 关闭定时任务
		room.Bait = false
		if timer := room.Timer; timer != nil {
			timer.Stop()
		}
		// 删除点歌歌单
		modules.MusicClear(room.ID)
	}
}

// HandleMember 处理用户相关的事件
func HandleMember(outputMsg chan<- string, store *models.Room, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventJoinQueue:
		// 有用户进入直播间
		for _, v := range textMsg.Queue {
			store.Count++ // 计数
			if username := v.Username; username != "" {
				text := fmt.Sprintf(models.TplWelcome, username)
				if store.Pinyin {
					// 如果注音功能开启了，发送注音消息
					text += fmt.Sprintf("\n注音：[%s]", utils.Pinyin(username))
				}
				outputMsg <- text
			} else if store.Count > 1 && store.Count%2 == 0 {
				// 屏蔽第一次匿名用户欢迎，减半欢迎匿名用户次数
				outputMsg <- models.TplWelcomeAnon
			}
		}
	case models.EventFollowed:
		// 有新关注
		if username := textMsg.User.Username; username != "" {
			outputMsg <- fmt.Sprintf(models.TplThankFollow, username)
		}
	}
}

// HandleGift 处理礼物相关的事件
func HandleGift(outputMsg chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventSend:
		// 有用户送礼物
		if username := textMsg.User.Username; username != "" {
			gift := textMsg.Gift
			outputMsg <- fmt.Sprintf(models.TplThankGift, username, gift.Number, gift.Name)
		}
	}
}

// HandleMessage 处理消息事件
func HandleMessage(outputMsg chan<- string, room *models.Room, textMsg models.FmTextMessage) {
	switch textMsg.Event {
	case models.EventNew:
		first := strings.Split(textMsg.Message, " ")[0]
		if first == "" {
			return
		}
		// 判断是否是命令，进行处理
		if cmdType := models.Command(first); cmdType >= 0 {
			handleCommand(outputMsg, room, cmdType, textMsg)
			return
		}
		// 判断是否是沟通请求，进行处理
		if first == fmt.Sprintf("@%s", config.Name()) {
			handleChat(outputMsg, room, textMsg)
			return
		}
		// 剩余的文本发送到关键词处理函数
		handleKeyword(outputMsg, room, textMsg)
	}
}

// handleCommand 处理消息中的命令，
// 简单的逻辑在本函数中处理，其余在 command.go 中处理
func handleCommand(outputMsg chan<- string, store *models.Room, cmdType int, textMsg models.FmTextMessage) {
	info, err := modules.RoomInfo(store.ID)
	if err != nil {
		zap.S().Error("获取直播间信息错误：", err)
		return
	}

	user := &textMsg.User // 当前发信的用户
	cmd := &command{
		Room:   store,
		User:   user,
		Role:   userRole(&info, user.UserID), // 获取当前发信用户的角色
		Output: outputMsg,
	}

	arr := strings.Fields(strings.TrimSpace(textMsg.Message))

	switch cmdType {
	case models.CmdRoomInfo:
		cmd.roomInfo(&info)
	case models.CmdCheckin:
		cmd.checkin(textMsg.User)
	case models.CmdCheckinRank:
		cmd.checkinRank()
	case models.CmdHoroscope:
		// check args
		if len(arr) == 2 {
			cmd.horoscopes(arr[1])
		}
	case models.CmdBaitSwitch:
		cmd.baitSwitch()
	case models.CmdWeather:
		// check args
		if len(arr) == 2 {
			cmd.weather(arr[1])
		}
	case models.CmdMusicAdd:
		// check args
		if len(arr) == 2 {
			cmd.musicAdd(arr[1])
		}
	case models.CmdMusicAll:
		cmd.musicAll()
	case models.CmdMusicPop:
		cmd.musicPop()
	case models.CmdPiaStart:
		// check args
		if len(arr) == 2 {
			id, err := strconv.Atoi(arr[1])
			if err != nil {
				return
			}
			cmd.piaStart(id)
		}
	case models.CmdPiaNext:
		if len(arr) == 1 {
			cmd.piaNext(1, false)
		}
		if len(arr) == 2 {
			dur, err := strconv.Atoi(arr[1])
			if err != nil {
				return
			}
			cmd.piaNext(dur, false)
		}
	case models.CmdPiaNextSafe:
		cmd.piaNextSafe()
	case models.CmdPiaRelocate:
		if len(arr) == 2 {
			idx, err := strconv.Atoi(arr[1])
			if err != nil {
				return
			}
			cmd.piaRelocate(idx)
		}
	case models.CmdPiaStop:
		cmd.piaStop()
	case models.CmdLove:
		outputMsg <- "❤️~"
	case models.CmdHelper:
		fallthrough
	default:
		outputMsg <- models.HelpText
	}
}

// handleChat 处理聊天请求
func handleChat(outputMsg chan<- string, store *models.Room, textMsg models.FmTextMessage) {
	outputMsg <- Chat(textMsg.User.Username)
}

// handleKeyword 处理关键词
func handleKeyword(outputMsg chan<- string, store *models.Room, textMsg models.FmTextMessage) {
	if strings.Contains(textMsg.Message, "emo") {
		keyEmotional(outputMsg, textMsg.User)
	}
}

// handleGame 处理游戏请求
func handleGame(outputMsg chan<- string, store *models.Room, textMsg *models.FmTextMessage) {

}
