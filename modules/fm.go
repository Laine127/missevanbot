package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"missevanbot/models"
)

const (
	Unfollow = 0
	Follow   = 1
)

// bot is used to store the basic information of the bot user.
var bot models.FmUser

// BotName return name of the bot.
func BotName() string {
	return bot.Username
}

// BotID return UID of the bot.
func BotID() int {
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
}

// UserInfo return the information of the bot user.
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

// RoomInfo return the information of the fm room.
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

// The BaseCookie get a new base cookie.
func BaseCookie() (string, error) {
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

// ChangeAttention change the attention state,
// tp = 0 is unfollow，tp = 1 is follow.
func ChangeAttention(uid, tp int) (ret []byte, err error) {
	_url := "https://www.missevan.com/person/ChangeAttention"

	data := []byte(fmt.Sprintf("attentionid=%d&type=%d", uid, tp))

	header := http.Header{}
	header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")

	ret, err = PostRequest(_url, header, data)
	return
}

// QueryUsername return the username queried by UID.
func QueryUsername(uid int) (string, error) {
	_url := fmt.Sprintf("https://www.missevan.com/%d/", uid)

	resp, err := http.Get(_url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	username := doc.Find("#t_u_n>a").Text()
	if username == "" {
		return "", errors.New("username empty")
	}

	return username, nil
}

// SendMessage send a private message to a user according to uid.
func SendMessage(uid int, content string) (ret []byte, err error) {
	_url := "https://www.missevan.com/mperson/sendmessage"

	data := []byte(fmt.Sprintf("user_id=%d&content=%s", uid, content))

	header := http.Header{}
	header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")

	ret, err = PostRequest(_url, header, data)
	return
}
