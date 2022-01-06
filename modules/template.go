package modules

import (
	"strings"
	"text/template"
)

const (
	TmplHelper         = "helper"
	TmplHelperPlaylist = "helper:playlist"
	TmplHelperPia      = "helper:pia"
	TmplHelperGame     = "helper:game"
	TmplUpdates        = "features"
	TmplMode           = "mode"
	TmplPlaylist       = "playlist"
	TmplRoomInfo       = "room"
	TmplRoomAdmin      = "admin"
	TmplStartUp        = "startup"
	TmplWelcome        = "welcome"
	TmplGift           = "gift"
	TmplNewAdmin       = "newadmin"
)

const (
	WordChat    = "chat"
	WordComfort = "comfort"
	WordGuess   = "guess"
	WordLuck    = "luck"
	WordReply   = "reply"
	WordWelcome = "welcome"
)

// The NewTemplate gets a text template via key from the Redis database,
// and use the data to generate the result text.
func NewTemplate(key string, data interface{}) (string, error) {
	temp := tmplString(key)
	text := new(strings.Builder)
	fn := template.FuncMap{"add": add}
	tmpl := template.Must(template.New(key).Funcs(fn).Parse(temp)) // hello.tmpl
	if err := tmpl.Execute(text, data); err != nil {
		return "", err
	}
	return text.String(), nil
}

func tmplString(key string) string {
	key = RedisPrefixTemplates + key

	c := rdb.Get(ctx, key)
	return c.Val()
}

func add(x, y int) int {
	return x + y
}

// The Word gets a random string from the Redis words set.
func Word(key string) string {
	key = RedisPrefixWords + key

	c := rdb.SRandMember(ctx, key)
	return c.Val()
}
