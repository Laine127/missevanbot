package models

type (
	// FmTextMessage represents the Websocket message from the live room.
	FmTextMessage struct {
		Type       string       `json:"type"`
		Event      string       `json:"event"`
		RoomID     int          `json:"room_id"`
		Message    string       `json:"message"`
		MessageID  string       `json:"msg_id"`
		User       FmUser       `json:"user"`
		Queue      []FmQueue    `json:"queue"`
		Gift       FmGift       `json:"gift"`
		Statistics FmStatistics `json:"statistics"`
	}

	// FmUser represents the information of a user.
	FmUser struct {
		IconUrl string `json:"iconurl"`
		Titles  []struct {
			Type  string `json:"type"`
			Name  string `json:"name"`
			Level int    `json:"level"`
		} `json:"titles"`
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
	}

	// FmQueue represents basic information of the user who is joining.
	FmQueue struct {
		Contribution int    `json:"contribution"`
		IconUrl      string `json:"iconurl"`
		UserId       int    `json:"user_id"`
		Username     string `json:"username"`
	}

	// FmGift represents the information of gift.
	FmGift struct {
		GiftID int    `json:"gift_id"`
		Name   string `json:"name"`
		Price  int    `json:"price"`
		Number int    `json:"num"`
	}

	// FmStatistics represents the statistics of the live room.
	FmStatistics struct {
		Accumulation int `json:"accumulation"`
		Online       int `json:"online"`
		Vip          int `json:"vip"`
		Score        int `json:"score"`
	}
)

// Event define the message events.
const (
	EventSend      = "send"       // gift send.
	EventNew       = "new"        // new message received.
	EventStatistic = "statistics" // statistics of the live room.
	EventJoin      = "join"       // connect to the live room channel.
	EventJoinQueue = "join_queue" // members join the live room.
	EventFollowed  = "followed"   // user followed the room creator.
	EventOpen      = "open"       // the live room opened.
	EventClose     = "close"      // the live room closed.
	EventNewRank   = "new_rank"   // the new rank information of the live room.
	EventLeave     = "leave"      // user leaved the live room.
)

// Type define the message types.
const (
	TypeRoom    = "room"
	TypeCreator = "creator"
	TypeGift    = "gift"
	TypeMessage = "message"
	TypeNotify  = "notify"
	TypeMember  = "member"
	TypeChannel = "channel"
)
