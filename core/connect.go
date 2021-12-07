package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"missevan-fm/models"
	"missevan-fm/modules"
)

var mu = &sync.Mutex{}

// Connect Websocket 连接处理
func Connect(inputMsg chan<- models.FmTextMessage, roomID int) {
	follow(roomID) // 关注主播

	dialer := new(websocket.Dialer)

	cookie, err := modules.ConnCookie()
	if err != nil {
		zap.S().Error("get the connection cookie failed: ", err)
		return
	}

	h := http.Header{}
	h.Add("Pragma", "no-cache")
	h.Add("Origin", "https://fm.missevan.com")
	h.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6,ja;q=0.5")
	h.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36 Edg/96.0.1054.34")
	h.Add("Cache-Control", "no-cache")
	h.Add("Cookie", cookie)

	conn, _, err := dialer.Dial(fmt.Sprintf("wss://im.missevan.com/ws?room_id=%d", roomID), h)
	if err != nil {
		zap.S().Error(err)
		return
	}
	defer conn.Close()

	joinMsg := fmt.Sprintf(`{"action":"join","uuid":"35e77342-30af-4b0b-a0eb-f80a826a68c7","type":"room","room_id":%d}`, roomID)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(joinMsg)); err != nil {
		zap.S().Error(err)
		return
	}

	go heart(conn) // 定时发送心跳，防止断开

	for {
		msgType, msgData, err := conn.ReadMessage()
		if err != nil {
			zap.S().Error(err)
			continue
		}

		zap.S().Debug(string(msgData))

		switch msgType {
		case websocket.TextMessage:
			if string(msgData) == "❤️" {
				continue // 过滤掉心跳消息
			}
			textMsg := models.FmTextMessage{}
			if err := json.Unmarshal(msgData, &textMsg); err != nil {
				zap.S().Errorf("%s 解析失败：%s", string(msgData), err)
				continue
			}
			if textMsg.User.UserID == modules.UserID() {
				continue // 过滤掉机器人自己的消息
			}
			inputMsg <- textMsg
		case websocket.BinaryMessage:
		case websocket.CloseMessage:
		case websocket.PingMessage:
		case websocket.PongMessage:
		default:
		}
	}
}

// heart 每隔一段时间就发送心跳
func heart(conn *websocket.Conn) {
	for {
		mu.Lock()
		_ = conn.WriteMessage(websocket.TextMessage, []byte("❤️"))
		mu.Unlock()
		time.Sleep(time.Second * 30) // 间隔时间 30s
	}
}

// follow 用于启动时关注主播
func follow(roomID int) {
	info, err := modules.RoomInfo(roomID)
	if err != nil {
		zap.S().Error("获取直播间信息错误：", err)
		return
	}

	if err := modules.UnfollowAll(); err != nil {
		zap.S().Error("取关操作出错：", err)
	}

	ret, err := modules.ChangeAttention(info.Creator.UserID, 1)
	if err != nil {
		zap.S().Error("关注用户出错：", err)
		return
	}
	zap.S().Debug(string(ret))
}
