package handlers

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"missevan-fm/config"
)

// PostRequest 发送 POST 请求，传递 data
func PostRequest(_url string, data []byte) (err error) {
	cookie := readCookie()
	if cookie == "" {
		return errors.New("cookie is empty")
	}

	client := new(http.Client)
	req, err := http.NewRequest("POST", _url, bytes.NewReader(data))
	if err != nil {
		return
	}

	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("cookie", cookie)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36 Edg/95.0.1020.53")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	return
}

// readCookie 读取当前目录下的 Cookie 文件，返回内容
func readCookie() string {
	file, err := os.Open(config.Cookie())
	if err != nil {
		log.Println("read cookie failed", err.Error())
		os.Exit(1)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("read cookie failed", err.Error())
		os.Exit(1)
	}
	return string(content)
}
