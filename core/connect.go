package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"missevanbot/models"
	"missevanbot/modules"
)

type connection struct {
	conn *websocket.Conn
	mu   *sync.Mutex
}

// Connect handle the Websocket connection,
// put the received message into input channel.
//
// TODO: Websocket connection may throw EOF error,
//  should find a new way to retry connection.
func Connect(ctx context.Context, cancel context.CancelFunc, input chan<- models.FmTextMessage, room *models.Room) {
	defer cancel() // cancel Connect, Match and Send goroutines.

	rid := room.ID
	// Follow the room creator.
	if err := follow(room); err != nil {
		zap.S().Warn(room.Log("follow user failed", err))
	}

	baseCookie, err := modules.BaseCookie()
	if err != nil {
		zap.S().Error(room.Log("get the base cookie failed", err))
		return
	}

	h := http.Header{}
	h.Add("Pragma", "no-cache")
	h.Add("Origin", "https://fm.missevan.com")
	h.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6,ja;q=0.5")
	h.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36 Edg/96.0.1054.34")
	h.Add("Cache-Control", "no-cache")
	h.Add("Cookie", baseCookie)

	dialer := new(websocket.Dialer)

retry:
	conn, resp, err := dialer.Dial(fmt.Sprintf("wss://im.missevan.com/ws?room_id=%d", rid), h)
	if err != nil {
		zap.S().Error(room.Log("dial failed", err))
		return
	}
	defer conn.Close()

	if resp.StatusCode != 101 {
		zap.S().Error(room.Log("websocket connect failed", nil))
		return
	}

	c := &connection{conn, &sync.Mutex{}}
	joinMsg := fmt.Sprintf(`{"action":"join","uuid":"35e77342-30af-4b0b-a0eb-f80a826a68c7","type":"room","room_id":%d}`, rid)
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(joinMsg)); err != nil {
		zap.S().Error(room.Log("write message failed", err))
		return
	}

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	go heart(ctx, c) // keep connected

	for {
		msgType, msgData, err := c.conn.ReadMessage()
		if err != nil {
			zap.S().Error(room.Log("read message failed", err))
			modules.MustPush(fmt.Sprint(modules.TitleError, room.Creator), err.Error())
			cancel() // cancel the heart goroutine.
			goto retry
		}

		zap.S().Debug(room.Log(string(msgData), err))

		switch msgType {
		case websocket.TextMessage:
			if string(msgData) == "❤️" {
				continue
			}
			textMsg := models.FmTextMessage{}
			if err := json.Unmarshal(msgData, &textMsg); err != nil {
				zap.S().Warn(room.Log(fmt.Sprintf("unmarshal failed (%s)", string(msgData)), err))
				continue
			}
			if textMsg.User.UserID == room.BotID() {
				continue // filter out messages sent by the bot.
			}
			input <- textMsg
		case websocket.BinaryMessage:
		case websocket.CloseMessage:
		case websocket.PingMessage:
		case websocket.PongMessage:
		default:
		}
	}
}

// heart send a heartbeat to the room Websocket connection every 30 seconds,
// also receive a WithCancel Context, when the `Connect` goroutine is down,
// return the heart function and stop the ticker.
func heart(ctx context.Context, c *connection) {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()
			_ = c.conn.WriteMessage(websocket.TextMessage, []byte("❤️"))
			c.mu.Unlock()
		}
	}
}

// follow get the ID of the room creator using room information api,
// and follow the creator according to the ID.
//
// If there is an error, write the logs and return.
func follow(room *models.Room) error {
	info, err := modules.RoomInfo(room)
	if err != nil {
		return err
	}

	ret, err := modules.ChangeAttention(room.BotCookie, info.Creator.UserID, modules.Follow)
	if err != nil {
		return err
	}

	if r := string(ret); r != `{"success":true,"info":"关注成功"}` {
		return errors.New(r)
	}
	return nil
}
