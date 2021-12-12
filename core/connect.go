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

type Connection struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

// Connect handle the Websocket connection,
// put the received message into inputMsg channel.
func Connect(inputMsg chan<- models.FmTextMessage, roomID int) {
	follow(roomID) // follow the room creator.

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

	c := new(Connection)
	c.conn, _, err = dialer.Dial(fmt.Sprintf("wss://im.missevan.com/ws?room_id=%d", roomID), h)
	if err != nil {
		zap.S().Error(err)
		return
	}
	defer c.conn.Close()

	joinMsg := fmt.Sprintf(`{"action":"join","uuid":"35e77342-30af-4b0b-a0eb-f80a826a68c7","type":"room","room_id":%d}`, roomID)
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(joinMsg)); err != nil {
		zap.S().Error(err)
		return
	}

	go heart(c) // keep connected

	for {
		msgType, msgData, err := c.conn.ReadMessage()
		if err != nil {
			zap.S().Error(err)
			continue
		}

		zap.S().Debug(string(msgData))

		switch msgType {
		case websocket.TextMessage:
			if string(msgData) == "❤️" {
				continue
			}
			textMsg := models.FmTextMessage{}
			if err := json.Unmarshal(msgData, &textMsg); err != nil {
				zap.S().Errorf("%s unmarshal failed: %s", string(msgData), err)
				continue
			}
			if textMsg.User.UserID == modules.UserID() {
				continue // filter out messages sent by the bot.
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

// heart send a heartbeat every 30 seconds.
func heart(c *Connection) {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		_ = c.conn.WriteMessage(websocket.TextMessage, []byte("❤️"))
		c.mu.Unlock()
	}
}

// follow follow the creator of the room.
func follow(roomID int) {
	info, err := modules.RoomInfo(roomID)
	if err != nil {
		zap.S().Error("fetch room info failed: ", err)
		return
	}

	ret, err := modules.ChangeAttention(info.Creator.UserID, modules.Follow)
	if err != nil {
		zap.S().Error("follow user failed: ", err)
		return
	}
	zap.S().Debug(string(ret))
}
