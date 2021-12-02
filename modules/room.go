package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// response 直播间响应字段
type response struct {
	Code int  `json:"code"`
	Info Info `json:"info"`
}

type Info struct {
	Creator Creator `json:"creator"`
	Room    room    `json:"room"`
}

// Creator 创建者响应字段
type Creator struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

type room struct {
	RoomID       int        `json:"room_id"`      // 直播间ID
	Name         string     `json:"name"`         // 直播间名
	Announcement string     `json:"announcement"` // 公告
	Members      Member     `json:"members"`      // 直播间成员
	Statistics   statistics `json:"statistics"`   // 统计数据
	Status       status     `json:"status"`       // 状态信息
}

type Member struct {
	Admin []Admin `json:"admin"` // 管理员
}

type statistics struct {
	Accumulation int `json:"accumulation"` // 累计人数
	Vip          int `json:"vip"`          // 贵宾数量
	Score        int `json:"score"`        // 分数
	Online       int `json:"online"`       // 在线
}

type status struct {
	Channel struct {
		Event    string `json:"event"`
		Platform string `json:"platform"`
		Time     int64  `json:"time"`
		Type     string `json:"type"`
	} `json:"channel"`
}

type Admin struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

// RoomInfo 获取直播间信息
func RoomInfo(roomID int) (info *Info, err error) {
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
	res := new(response)
	if err = json.Unmarshal(body, res); err != nil {
		return
	}
	return &res.Info, nil
}
