package module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Room 直播间响应字段
type Room struct {
	Code int `json:"code"`
	Info struct {
		Creator Creator `json:"creator"`
	} `json:"info"`
}

// Creator 创建者响应字段
type Creator struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

// RoomInfo 获取直播间信息
func RoomInfo(roomID int) *Room {
	_url := fmt.Sprintf("https://fm.missevan.com/api/v2/live/%d", roomID)
	resp, err := http.Get(_url)
	if err != nil {
		log.Println("获取直播间信息错误：", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取直播间信息错误：", err)
		return nil
	}
	room := new(Room)
	if err := json.Unmarshal(body, room); err != nil {
		log.Println("解析直播间信息响应错误：", err)
		return nil
	}
	return room
}
