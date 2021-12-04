package modules

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"missevan-fm/config"
)

// PostRequest 发送 POST 请求，传递 data
func PostRequest(_url string, header http.Header, data []byte) (body []byte, err error) {
	cookie := config.Cookie()
	if cookie == "" {
		err = errors.New("cookie is empty")
		return
	}

	client := new(http.Client)
	req, err := http.NewRequest("POST", _url, bytes.NewReader(data))
	if err != nil {
		return
	}

	req.Header = header
	req.Header.Set("origin", "https://www.missevan.com")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36 Edg/95.0.1020.53")
	req.Header.Set("cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	return
}
