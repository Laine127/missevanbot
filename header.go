package main

import "net/http"

func Header() *http.Header {
	header := http.Header{}
	header.Set("Host", "im.missevan.com")
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:94.0) Gecko/20100101 Firefox/94.0")
	header.Set("Accept", "*/*")
	header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	header.Set("Accept-Encoding", "gzip, deflate, br")
	// header.Set("Sec-WebSocket-Version","13")
	// header.Set("Upgrade","websocket")
	header.Set("Origin", "https://fm.missevan.com")
	// header.Set("Sec-WebSocket-Extensions","permessage-deflate")
	// header.Set("Sec-WebSocket-Key","6anGJ9ZtrqfGmuWnakoFDw==")
	// header.Set("Connection","keep-alive, Upgrade")
	header.Set("Cookie", "FM_SESS=20211123|8jvxga0hfxz8uqwelo3pniyov; FM_SESS.sig=Crk9p_L0eW6YKtwBkFN0viuR1EU")
	header.Set("Sec-Fetch-Dest", "websocket")
	header.Set("Sec-Fetch-Mode", "websocket")
	header.Set("Sec-Fetch-Site", "same-site")
	header.Set("Pragma", "no-cache")
	header.Set("Cache-Control", "no-cache")

	return &header
}
