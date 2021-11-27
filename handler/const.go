package handler

// FmTextMessage 直播间Websocket消息主体
type FmTextMessage struct {
	Type      string `json:"type"`
	Event     string `json:"event"`
	RoomID    int    `json:"room_id"`
	Message   string `json:"message"`
	MessageID string `json:"msg_id"`

	User       FmUser       `json:"user"`
	Queue      []FmQueue    `json:"queue"`
	Gift       FmGift       `json:"gift"`
	Statistics FmStatistics `json:"statistics"`
}

// FmUser 直播间Websocket用户消息体
type FmUser struct {
	IconUrl string `json:"iconurl"`
	Titles  []struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		Level int    `json:"level"`
	} `json:"titles"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

// FmQueue 直播间Websocket观众队列消息体
type FmQueue struct {
	Contribution int    `json:"contribution"`
	IconUrl      string `json:"iconurl"`
	UserId       int    `json:"user_id"`
	Username     string `json:"username"`
}

// FmGift 直播间Websocket礼物消息体
type FmGift struct {
	GiftID int    `json:"gift_id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Number int    `json:"num"`
}

// FmStatistics 直播间Websocket静态信息消息体
type FmStatistics struct {
	Accumulation int `json:"accumulation"`
	Online       int `json:"online"`
	Vip          int `json:"vip"`
	Score        int `json:"score"`
}

// Event 定义事件
const (
	EventSend      = "send"       // 发送
	EventNew       = "new"        // 普通新消息
	EventStatistic = "statistics" // 直播间信息
	EventJoin      = "join"       // 主动连接直播间
	EventJoinQueue = "join_queue" // 用户进入直播间
	EventFollowed  = "followed"   // 用户关注了主播
	EventOpen      = "open"       // 直播间开启
	EventClose     = "close"      // 直播间关闭
	EventNewRank   = "new_rank"   // 直播间排行
	EventLeave     = "leave"      // 用户离开直播间
)

// Type 定义消息类型
const (
	TypeRoom    = "room"    // 直播间相关
	TypeCreator = "creator" // 创作者
	TypeGift    = "gift"    // 礼物
	TypeMessage = "message" // 消息
	TypeNotify  = "notify"  // 公屏提示
	TypeMember  = "member"  // 成员/用户
	TypeChannel = "channel" // 连接隧道
)
