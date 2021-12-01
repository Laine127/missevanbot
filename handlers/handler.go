package handlers

import (
	"fmt"
	"strings"

	"github.com/mozillazg/go-pinyin"
	"go.uber.org/zap"
	"missevan-fm/config"
	"missevan-fm/models"
	"missevan-fm/modules"
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
		room.Online = textMsg.Statistics.Online // 更新在线人数
	case models.EventOpen:
		outputMsg <- models.TplBotStart
		// 通知推送
		text := fmt.Sprintf("%s 开播啦~", info.Creator.Username)
		if err := modules.Push(modules.TitleOpen, text); err != nil {
			zap.S().Error("Bark 推送失败", err)
		}
	case models.EventClose:
		// 通知推送
		text := fmt.Sprintf("%s 下播啦~", info.Creator.Username)
		if err := modules.Push(modules.TitleClose, text); err != nil {
			zap.S().Error("Bark 推送失败", err)
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
			var text string

			store.Count++

			if username := v.Username; username != "" {
				text = fmt.Sprintf(models.TplWelcome, username)
				if store.Pinyin {
					// 如果注音功能开启了，发送注音消息
					py := pinyin.NewArgs()
					py.Style = pinyin.Tone
					if arr := pinyin.Pinyin(username, py); len(arr) > 0 {
						text += fmt.Sprintf("\n注音：%s", arr)
					}
				}
			} else if store.Count > 1 && store.Count%2 == 0 {
				// 屏蔽第一次匿名用户欢迎，减半欢迎匿名用户次数
				text = models.TplWelcomeAnon
			}
			outputMsg <- text
		}
	case models.EventFollowed:
		// 有新关注
		if username := textMsg.User.Username; username != "" {
			outputMsg <- fmt.Sprintf(models.TplThankFollow, username)
		}
	}
	return
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
		// 判断是否是命令，进行处理
		arr := strings.Split(textMsg.Message, " ")
		if cmdType := models.Command(arr[0]); cmdType >= 0 {
			handleCommand(outputMsg, room, cmdType, textMsg)
			return
		}
		// 判断是否是沟通请求，进行处理
		if arr[0] == fmt.Sprintf("@%s", config.Name()) {
			handleChat(outputMsg, room, textMsg)
			return
		}
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
		Role:   userRole(info, user.UserID), // 获取当前发信用户的角色
		Output: outputMsg,
	}

	arr := strings.Split(textMsg.Message, " ")

	switch cmdType {
	case models.CmdInfo:
		cmd.info(info)
	case models.CmdSign:
		cmd.sign(textMsg.User)
	case models.CmdRank:
		cmd.rank()
	case models.CmdBait:
		cmd.bait()
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

// handleGame 处理游戏请求
func handleGame(outputMsg chan<- string, store *models.Room, textMsg *models.FmTextMessage) {

}
