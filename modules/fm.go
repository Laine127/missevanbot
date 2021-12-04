package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"missevan-fm/models"
)

// RoomInfo 获取直播间信息
func RoomInfo(roomID int) (info models.FmInfo, err error) {
	_url := fmt.Sprintf("https://fm.missevan.com/api/v2/live/%d", roomID)
	resp, err := http.Get(_url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	res := new(models.FmRoom)
	if err = json.Unmarshal(body, res); err != nil {
		return
	}
	return res.Info, nil
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
