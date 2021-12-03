package thirdparty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Level 运势递增变好
const (
	_ = iota
	Level1
	Level2
	Level3
	Level4
	Level5
)

var StarList = map[string]struct{}{
	"白羊座": {},
	"金牛座": {},
	"双子座": {},
	"巨蟹座": {},
	"狮子座": {},
	"处女座": {},
	"天秤座": {},
	"天蝎座": {},
	"射手座": {},
	"摩羯座": {},
	"水瓶座": {},
	"双鱼座": {},
}

type Fortune struct {
	Content string `json:"content"`
	Date    string `json:"date"`
	Score   string `json:"score"`
}

// Zodiac 获取星座运势
func Zodiac(star string, level int) (fort *Fortune, err error) {
	_url := fmt.Sprintf("https://datamuse.guokr.com/api/front/common/muse/constellation/v1/fortune?constellation=%s&level=%d", star, level)
	resp, err := http.Get(_url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fort = new(Fortune)
	err = json.Unmarshal(body, fort)
	return
}
