package handler

import (
	"fmt"
	"strings"

	"missevan-fm/config"
)

// handleMessage 处理消息事件
func handleMessage(room *config.RoomConfig, textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventNew:
		if cmdType, ok := _cmdMap[textMsg.Message]; ok {
			// 检查是否是命令，进行处理
			handleCommand(room, cmdType, textMsg)
			return
		}
		if strings.HasPrefix(textMsg.Message, fmt.Sprintf("@%s", config.Conf.Name)) {
			// 检查是否是聊天请求，进行处理
			handleChat(room, textMsg)
		}
	}
}
