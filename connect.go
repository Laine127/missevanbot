package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var mu = &sync.Mutex{}

// header 返回特定的 Websocket 首部
func header() *http.Header {
	h := http.Header{}
	h.Set("Host", "im.missevan.com")
	h.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:94.0) Gecko/20100101 Firefox/94.0")
	h.Set("Accept", "*/*")
	h.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	h.Set("Accept-Encoding", "gzip, deflate, br")
	// h.Set("Sec-WebSocket-Version","13")
	// h.Set("Upgrade","websocket")
	h.Set("Origin", "https://fm.missevan.com")
	// h.Set("Sec-WebSocket-Extensions","permessage-deflate")
	// h.Set("Sec-WebSocket-Key","6anGJ9ZtrqfGmuWnakoFDw==")
	// h.Set("Connection","keep-alive, Upgrade")
	h.Set("Cookie", "FM_SESS=20211123|8jvxga0hfxz8uqwelo3pniyov; FM_SESS.sig=Crk9p_L0eW6YKtwBkFN0viuR1EU")
	h.Set("Sec-Fetch-Dest", "websocket")
	h.Set("Sec-Fetch-Mode", "websocket")
	h.Set("Sec-Fetch-Site", "same-site")
	h.Set("Pragma", "no-cache")
	h.Set("Cache-Control", "no-cache")
	return &h
}

// connect Websocket 连接处理
func connect(roomID int) {
	dialer := new(websocket.Dialer)

	conn, resp, err := dialer.Dial(fmt.Sprintf("wss://im.missevan.com/ws?room_id=%d", roomID), *header())
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	if resp.StatusCode != 101 {
		// 101 响应成功，非 101 则失败
		log.Println("请求失败")
		return
	}

	joinMsg := fmt.Sprintf(`{"action":"join","uuid":"35e77342-30af-4b0b-a0eb-f80a826a68c7","type":"room","room_id":%d}`, roomID)

	if err := conn.WriteMessage(websocket.TextMessage, []byte(joinMsg)); err != nil {
		log.Println(err)
		return
	}

	go heart(conn) // 定时发送心跳，防止断开

	for {
		msgType, msgData, err := conn.ReadMessage()
		if nil != err {
			log.Println(err)
			break
		}

		switch msgType {
		case websocket.TextMessage:
			// 处理文本消息
			handleTextMessage(roomID, string(msgData))
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
