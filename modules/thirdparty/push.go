package thirdparty

import (
	"fmt"
	"net/http"
)

// BarkPush 通过 Bark 推送
func BarkPush(token, title, msg string) (err error) {
	_url := fmt.Sprintf("https://api.day.app/%s/%s/%s", token, title, msg)

	resp, err := http.Get(_url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return
}
