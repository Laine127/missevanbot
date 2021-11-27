package message

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"missevan-fm/bot"
)

type message struct {
	RoomID    int    `json:"room_id"`
	Message   string `json:"message"`
	MessageID string `json:"msg_id"`
}

// MustSend 发送一条消息，输出错误日志并忽略
func MustSend(roomID int, text string) {
	if err := Send(roomID, text); err != nil {
		log.Println("发送消息出错：", err)
		return
	}
}

// Send 发送一条消息，返回错误
func Send(roomID int, text string) (err error) {
	_url := "https://fm.missevan.com/api/chatroom/message/send"

	cookie := readCookie()
	if cookie == "" {
		return errors.New("cookie is empty")
	}

	data, _ := json.Marshal(message{
		RoomID:    roomID,
		Message:   text,
		MessageID: messageID(),
	})

	client := new(http.Client)
	req, err := http.NewRequest("POST", _url, bytes.NewReader(data))
	if err != nil {
		return
	}

	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("cookie", cookie)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36 Edg/95.0.1020.53")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)

	return
}

// readCookie 读取当前目录下的 Cookie 文件，返回内容
func readCookie() string {
	conf := bot.Conf.Cookie
	file, err := os.Open(conf)
	if err != nil {
		log.Println("读取Cookie失败啦：", err)
		os.Exit(1)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("读取Cookie失败啦：", err)
		os.Exit(1)
	}
	return string(content)
}

// messageID 用于生成消息唯一ID
func messageID() string {
	u := uuid.NewString()
	return "3" + u[1:]
}
