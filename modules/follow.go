package modules

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"missevan-fm/config"
)

// MustFollow 关注动作，不返回错误
func MustFollow(uid int) {
	if err := Follow(uid); err != nil {
		log.Println("关注用户出错", err.Error())
		return
	}
}

// Follow 关注动作，返回错误
func Follow(uid int) (err error) {
	_url := "https://www.missevan.com/person/ChangeAttention"

	cookie := config.Cookie()
	if cookie == "" {
		return errors.New("cookie is empty")
	}

	data := []byte(fmt.Sprintf("attentionid=%d&type=1", uid))

	client := new(http.Client)
	req, err := http.NewRequest("POST", _url, bytes.NewReader(data))
	if err != nil {
		return
	}

	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("cookie", cookie)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36 Edg/96.0.1054.34")
	req.Header.Set("origin", "https://www.missevan.com")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)

	log.Println("关注用户", fmt.Sprintf("%d", uid))

	return
}
