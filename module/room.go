package module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

type Admin struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

// RoomInfo 获取直播间信息
func RoomInfo(roomID int) *Info {
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
	res := new(response)
	if err := json.Unmarshal(body, res); err != nil {
		log.Println("解析直播间信息响应错误：", err)
		return nil
	}
	return &res.Info
}
