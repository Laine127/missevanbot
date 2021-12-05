package modules

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"missevan-fm/models"
)

var bot models.FmUser // 存储机器人用户名

// Name return name of the bot.
func Name() string {
	return bot.Username
}

// UserID return UID of the bot.
func UserID() int {
	return bot.UserID
}

// InitBot initialize information of the bot.
func InitBot() {
	// get name of the bot.
	user, err := UserInfo()
	if err != nil {
		panic(fmt.Errorf("got bot name failed: %s", err))
	}
	bot = user
	log.Println(user)
}

// UserInfo 获取机器人信息
func UserInfo() (user models.FmUser, err error) {
	_url := "https://fm.missevan.com/api/user/info"

	body, err := GetRequest(_url, nil)
	if err != nil {
		return
	}

	info := struct {
		Code int `json:"code"`
		Info struct {
			User models.FmUser `json:"user"`
		} `json:"info"`
	}{}

	err = json.Unmarshal(body, &info)
	user = info.Info.User
	return
}

// RoomInfo 获取直播间信息
func RoomInfo(roomID int) (info models.FmInfo, err error) {
	_url := fmt.Sprintf("https://fm.missevan.com/api/v2/live/%d", roomID)

	body, err := GetRequest(_url, nil)
	if err != nil {
		return
	}

	room := models.FmRoom{}
	if err = json.Unmarshal(body, &room); err != nil {
		return
	}

	info = room.Info
	return
}

// ConnCookie 获取初始连接的 Cookie 字符串
func ConnCookie() (string, error) {
	_url := "https://fm.missevan.com/api/user/info"

	resp, err := http.Get(_url)
	if err != nil {
		return "", err
	}
	cookie := strings.Builder{}
	for _, v := range resp.Header.Values("set-cookie") {
		cookie.WriteString(v)
	}
	return cookie.String(), nil
}

// Follow 关注动作，返回错误
func Follow(uid int) (ret []byte, err error) {
	_url := "https://www.missevan.com/person/ChangeAttention"

	data := []byte(fmt.Sprintf("attentionid=%d&type=1", uid))

	header := http.Header{}
	header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")

	ret, err = PostRequest(_url, header, data)
	return
}
