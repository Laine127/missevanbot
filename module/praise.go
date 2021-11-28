package module

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"missevan-fm/bot"
)

// Praise 彩虹屁模块
func Praise(room *bot.RoomConfig, timer *time.Timer) {
	for {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randNumber := r.Intn(room.RainbowMaxInterval)
		timer.Reset(time.Duration(randNumber+1) * time.Minute)
		<-timer.C

		if text := comment(); text != "" {
			MustSend(room.ID, text)
		}
	}
}

// comment 获取彩虹屁接口内容
func comment() (s string) {
	resp, err := http.Get("https://api.muxiaoguo.cn/api/caihongpi")
	if err != nil {
		ll.Print("请求彩虹屁接口失败", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ll.Print("读取彩虹屁接口失败", err.Error())
		return
	}
	// 彩虹屁接口 JSON 结构体
	c := new(struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Comment string `json:"comment"`
		} `json:"data"`
	})
	if err := json.Unmarshal(body, c); err != nil {
		ll.Print("解析彩虹屁接口失败", err.Error())
	}
	if c.Code != "200" || c.Msg != "success" {
		ll.Print("请求彩虹屁接口错误", c.Code)
		return
	}
	return c.Data.Comment
}
