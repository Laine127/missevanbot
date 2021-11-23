package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"missevan-fm/module"
)

var mu = &sync.Mutex{}

func connect(roomID int) {
	dialer := websocket.Dialer{}

	connect, resp, err := dialer.Dial(fmt.Sprintf("wss://im.missevan.com/ws?room_id=%d", roomID), *Header())
	if err != nil {
		log.Println(err)
		return
	}
	defer connect.Close()

	if resp.StatusCode != 101 {
		log.Println("请求失败")
		return
	}

	joinMsg := fmt.Sprintf(`{"action":"join","uuid":"35e77342-30af-4b0b-a0eb-f80a826a68c7","type":"room","room_id":%d}`, roomID)

	if err := connect.WriteMessage(websocket.TextMessage, []byte(joinMsg)); err != nil {
		log.Println(err)
		return
	}

	go heart(connect)
	// read from server
	for {
		msgType, msgData, err := connect.ReadMessage()
		if nil != err {
			log.Println(err)
			break
		}

		switch msgType {
		case websocket.TextMessage:
			handleTextMessage(roomID, string(msgData))
		case websocket.BinaryMessage:
			// fmt.Println(msgData)
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
		_ = conn.WriteMessage(websocket.TextMessage, []byte(`❤️`))
		mu.Unlock()
		time.Sleep(time.Second * 30) // 间隔时间 30s
	}
}

func handleTextMessage(roomID int, msg string) {
	fmt.Println(msg)
	if msg == "❤️" {
		// 过滤掉心跳消息
		return
	}
	data := new(FmTextMessage)
	if err := json.Unmarshal([]byte(msg), data); err != nil {
		log.Println("解析失败了：", err)
		return
	}
	if data.Event == "join_queue" {
		// 加入房间
		for _, v := range data.Queue {
			username := v.Username
			if username == "" {
				continue
			}
			module.SendMessage(roomID, fmt.Sprintf("欢迎%s进入直播间~", username))
		}
	}
}
