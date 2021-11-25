package handler

import "missevan-fm/config"

// handleMessage 处理消息事件
func handleMessage(room *config.RoomConfig, textMsg *FmTextMessage) {
	switch textMsg.Event {
	case EventNew:
		if cmdType, ok := _cmdMap[textMsg.Message]; ok {
			// 命令处理
			handleCommand(room, cmdType, textMsg)
		}
	}
}
