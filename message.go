package main

type FmTextMessage struct {
	Event      string       `json:"event"`
	Message    string       `json:"message"`
	MessageID  string       `json:"msg_id"`
	RoomID     int          `json:"room_id"`
	Type       string       `json:"type"`
	User       FmUser       `json:"user"`
	Queue      []FmQueue    `json:"queue"`
	Gift       FmGift       `json:"gift"`
	Statistics FmStatistics `json:"statistics"`
}

type FmUser struct {
	IconUrl string `json:"iconurl"`
	Titles  []struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		Level int    `json:"level"`
	} `json:"titles"`
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

type FmQueue struct {
	Contribution int    `json:"contribution"`
	IconUrl      string `json:"iconurl"`
	UserId       int    `json:"user_id"`
	Username     string `json:"username"`
}

type FmGift struct {
	GiftID int    `json:"gift_id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Number int    `json:"num"`
}

type FmStatistics struct {
	Accumulation int `json:"accumulation"`
	Online       int `json:"online"`
	Vip          int `json:"vip"`
	Score        int `json:"score"`
}

const (
	EventSend      = "send"       // 发送
	EventNew       = "new"        // 普通新消息
	EventStatistic = "statistics" // 直播间信息
	EventJoin      = "join_queue" // 用户进入直播间
	EventFollowed  = "followed"   // 用户关注了主播
)

const (
	TypeRoom    = "room"    // 直播间相关
	TypeCreator = "creator" // 创作者
	TypeGift    = "gift"    // 礼物
	TypeMessage = "message" // 消息
	TypeNotify  = "notify"  // 公屏提示
	TypeMember  = "member"  // 成员/用户
)
