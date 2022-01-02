package modules

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const (
	ReqGet  = 1
	ReqPost = 2
)

type Request struct {
	URL    string
	Header http.Header
	Cookie string
	Data   []byte
}

func NewRequest(url string, hdr http.Header, c string, data []byte) Request {
	return Request{
		URL:    url,
		Header: hdr,
		Cookie: c,
		Data:   data,
	}
}

// Post does the POST request with the data.
func (r Request) Post() (body []byte, err error) {
	return r.request(ReqPost)
}

// Get does the GET request.
func (r Request) Get() (body []byte, err error) {
	return r.request(ReqGet)
}

// request does the GET or the POST request according to reqType,
// first, read the cookie of the bot user.
func (r Request) request(reqType int) (body []byte, err error) {
	client, req := new(http.Client), new(http.Request)

	switch reqType {
	case ReqGet:
		req, err = http.NewRequest("GET", r.URL, nil)
	case ReqPost:
		req, err = http.NewRequest("POST", r.URL, bytes.NewReader(r.Data))
	}
	if err != nil {
		return
	}

	req.Header = http.Header{}
	if r.Header != nil {
		req.Header = r.Header
	}
	req.Header.Set("origin", "https://www.missevan.com")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36 Edg/95.0.1020.53")
	req.Header.Set("cookie", r.Cookie)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	return
}
