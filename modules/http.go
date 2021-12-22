package modules

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"missevanbot/config"
)

const (
	ReqGet  = 1
	ReqPost = 2
)

// PostRequest does the POST request with the data.
func PostRequest(_url string, header http.Header, data []byte) (body []byte, err error) {
	return request(_url, header, data, ReqPost)
}

// GetRequest does the GET request.
func GetRequest(_url string, header http.Header) (body []byte, err error) {
	return request(_url, header, nil, ReqGet)
}

// request does the GET or the POST request according to reqType,
// first, read the cookie of the bot user.
func request(_url string, header http.Header, data []byte, reqType int) (body []byte, err error) {
	cookie := config.Cookie()
	if cookie == "" {
		err = errors.New("cookie is empty")
		return
	}

	client, req := new(http.Client), new(http.Request)

	switch reqType {
	case ReqGet:
		req, err = http.NewRequest("GET", _url, nil)
	case ReqPost:
		req, err = http.NewRequest("POST", _url, bytes.NewReader(data))
	}
	if err != nil {
		return
	}

	req.Header = http.Header{}
	if header != nil {
		req.Header = header
	}
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
