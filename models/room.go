package models

import (
	"container/list"
	"fmt"
	"time"

	"missevanbot/config"
)

// Room 直播间实例结构体
type Room struct {
	*config.RoomConfig              // 当前房间的配置
	Count              int          // 统计进入的数量
	Online             int          // 记录当前直播间在线人数
	Playlist           *list.List   //
	PiaList            []string     // 戏文
	PiaIndex           int          // 位置
	Ticker             *time.Ticker //
	TickerCount        int          //
	Gamer              Gamer        // 存储游戏状态
}

func NewRoom(roomConf *config.RoomConfig) *Room {
	return &Room{RoomConfig: roomConf, Playlist: list.New()}
}

func (room *Room) Log(str string, err interface{}) string {
	if err == nil {
		return fmt.Sprintf("[%s][%d] %s", room.Creator, room.ID, str)
	}
	return fmt.Sprintf("[%s][%d] %s: %s", room.Creator, room.ID, str, err)
}

type (
	FmRoom struct {
		Code int    `json:"code"`
		Info FmInfo `json:"info"`
	}

	FmInfo struct {
		Creator fmCreator `json:"creator"`
		Room    fmRoom    `json:"room"`
	}

	fmCreator struct {
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
	}

	fmRoom struct {
		RoomID       int          `json:"room_id"`      // 直播间ID
		Name         string       `json:"name"`         // 直播间名
		Announcement string       `json:"announcement"` // 公告
		Members      fmMembers    `json:"members"`      // 直播间成员
		Statistics   fmStatistics `json:"statistics"`   // 统计数据
		Status       fmStatus     `json:"status"`       // 状态信息
	}

	fmMembers struct {
		Admin []fmAdmin `json:"admin"` // 管理员
	}

	fmAdmin struct {
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
	}

	fmStatistics struct {
		Accumulation   int `json:"accumulation"`    // 累计人数
		Vip            int `json:"vip"`             // 贵宾数量
		Score          int `json:"score"`           // 分数
		Online         int `json:"online"`          // 在线
		AttentionCount int `json:"attention_count"` // 关注数
	}

	fmStatus struct {
		Open    int       `json:"open"`
		Channel fmChannel `json:"channel"`
	}

	fmChannel struct {
		Event    string `json:"event"`
		Platform string `json:"platform"`
		Time     int64  `json:"time"`
		Type     string `json:"type"`
	}
)
