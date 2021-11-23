package main

type FmTextMessage struct {
	Event     string `json:"event"`
	Message   string `json:"message"`
	MessageID string `json:"msg_id"`
	RoomID    int    `json:"room_id"`
	Type      string `json:"type"`
	User      FmUser `json:"user"`

	Queue []FmQueue `json:"queue"`
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
