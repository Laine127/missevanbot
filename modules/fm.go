package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"missevanbot/config"
	"missevanbot/models"
)

const (
	Unfollow = 0
	Follow   = 1
)

// InitBot queries information of the bot.
func InitBot(room *models.Room) (err error) {
	c := cookie(room.ID)
	if c == "" {
		// Make sure the cookie is not empty.
		err = errors.New("cookie empty")
		return
	}
	room.BotCookie = c
	user, err := UserInfo(c)
	if err != nil {
		return
	}
	room.BotUser = user
	return
}

// UserInfo return the information of the bot user.
func UserInfo(c string) (user models.FmUser, err error) {
	_url := "https://fm.missevan.com/api/user/info"

	req := NewRequest(_url, nil, c, nil)
	body, err := req.Get()
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
func RoomInfo(room *models.Room) (info models.FmInfo, err error) {
	_url := fmt.Sprintf("https://fm.missevan.com/api/v2/live/%d", room.ID)

	req := NewRequest(_url, nil, room.BotCookie, nil)
	body, err := req.Get()
	if err != nil {
		return
	}

	fmRoom := models.FmRoom{}
	if err = json.Unmarshal(body, &fmRoom); err != nil {
		return
	}
	info = fmRoom.Info
	return
}

// The BaseCookie get a new base cookie.
func BaseCookie() (string, error) {
	_url := "https://fm.missevan.com/api/user/info"

	resp, err := http.Get(_url)
	if err != nil {
		return "", err
	}
	bc := strings.Builder{}
	for _, v := range resp.Header.Values("set-cookie") {
		bc.WriteString(v)
	}
	return bc.String(), nil
}

// ChangeAttention change the attention state,
// tp = 0 is unfollow，tp = 1 is follow.
func ChangeAttention(c string, uid, tp int) (ret []byte, err error) {
	_url := "https://www.missevan.com/person/ChangeAttention"

	data := []byte(fmt.Sprintf("attentionid=%d&type=%d", uid, tp))

	header := http.Header{}
	header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")

	req := NewRequest(_url, header, c, data)
	ret, err = req.Post()
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

	name := doc.Find("#t_u_n>a").Text()
	if name == "" {
		return "", errors.New("username empty")
	}

	return name, nil
}

// SendMessage send a private message to a user according to uid.
func SendMessage(c string, uid int, content string) (ret []byte, err error) {
	_url := "https://www.missevan.com/mperson/sendmessage"

	data := []byte(fmt.Sprintf("user_id=%d&content=%s", uid, content))

	header := http.Header{}
	header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")

	req := NewRequest(_url, header, c, data)
	ret, err = req.Post()
	return
}

const DefaultCookie = "0"

func cookie(rid int) string {
	rdb := config.RDB
	key := config.RedisPrefixCookies
	field := strconv.Itoa(rid)
	if !rdb.HExists(ctx, key, field).Val() {
		field = DefaultCookie
	}
	return rdb.HGet(ctx, key, field).Val()
}
