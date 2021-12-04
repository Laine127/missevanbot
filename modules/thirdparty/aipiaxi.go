package thirdparty

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Fetch 获取戏文
func Fetch(id int) (roleList []string, textList []string, err error) {
	_url := fmt.Sprintf("https://aipiaxi.com/Index/post/id/%d", id)

	client := http.Client{}
	req, err := http.NewRequest("GET", _url, nil)
	req.Header.Add("Cookie", "XJUID=MTYzODQ2MTIxNjgxNjAuMDQ4OTcxNDAwNjcyODkzNDA2")
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = errors.New("status error")
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	// 获取角色列表
	roleList = make([]string, 0)
	doc.Find(".role-section>div>*").Each(func(i int, s *goquery.Selection) {
		roleList = append(roleList, s.Text())
	})

	// 获取文本段列表
	textList = make([]string, 0)
	doc.Find("#content>*").Each(func(i int, s *goquery.Selection) {
		content := s.Find("*").Text()
		text := s.Text()
		text = strings.TrimSpace(text)
		content = strings.TrimSpace(content)
		if s.Is("hr,br") || text == " " {
			return
		} else if text != "" {
			textList = append(textList, s.Text())
		} else if content != "" {
			textList = append(textList, content)
		}
	})
	return
}
