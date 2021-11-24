package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"missevan-fm/config"
)

type message struct {
	RoomID    int    `json:"room_id"`
	Message   string `json:"message"`
	MessageID string `json:"msg_id"`
}

// SendMessage 发送一条消息 Post
func SendMessage(roomID int, msg string) {
	if msg == "" {
		log.Println("发送的消息为空。。。")
		return
	}
	_url := "https://fm.missevan.com/api/chatroom/message/send"

	cookie := readCookie()

	if cookie == "" {
		log.Println("Cookie是空的。。。")
		return
	}

	d := message{
		RoomID:    roomID,
		Message:   msg,
		MessageID: messageID(),
	}
	data, _ := json.Marshal(d)

	fmt.Println(string(data))

	client := new(http.Client)
	req, err := http.NewRequest("POST", _url, bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("cookie", cookie)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36 Edg/95.0.1020.53")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(body))
}

// readCookie 读取当前目录下的 `.cookie` 文件，返回内容
func readCookie() string {
	conf := config.Conf.Cookie
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
