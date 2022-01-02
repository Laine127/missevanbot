package modules

import (
	"strings"
	"text/template"

	"missevanbot/config"
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
)

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
	rdb := config.RDB
	key = config.RedisPrefixTemplates + key
	c := rdb.Get(ctx, key)
	return c.Val()
}

func add(x, y int) int {
	return x + y
}
