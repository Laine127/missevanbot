package thirdparty

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Fetch 获取戏文
func Fetch(id int) ([]string, error) {
	_url := fmt.Sprintf("https://aipiaxi.com/Index/post/id/%d", id)

	client := http.Client{}
	req, err := http.NewRequest("GET", _url, nil)
	req.Header.Add("Cookie", "XJUID=MTYzODQ2MTIxNjgxNjAuMDQ4OTcxNDAwNjcyODkzNDA2")
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("status error")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	textList := make([]string, 0)

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
	return textList, nil
}
