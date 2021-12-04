package thirdparty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// WeatherText 获取天气
func WeatherText(city string) (text string, err error) {
	resp, err := http.Get(fmt.Sprintf("https://api.muxiaoguo.cn/api/tianqi?city=%s&type=1", city))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// 天气接口 JSON 结构体
	c := struct {
		Code interface{} `json:"code"`
		Msg  string      `json:"msg"`
		Data struct {
			CityName string `json:"cityname"` // 城市名
			Weather  string `json:"weather"`  // 天气
			Temp     string `json:"temp"`     // 温度
			SD       string `json:"SD"`       // 相对湿度
		} `json:"data"`
	}{}

	if err = json.Unmarshal(body, &c); err != nil {
		return
	}

	if v, ok := c.Code.(int); ok && v == -1 {
		return c.Msg, nil // request failed
	}

	return fmt.Sprintf("%s 今日天气：\n[天气] %s\n[温度] %s\n[相对湿度] %s",
		c.Data.CityName,
		c.Data.Weather,
		c.Data.Temp,
		c.Data.SD,
	), nil
}

// PoemText 获取一句诗词
func PoemText() (text string, err error) {
	_url := "https://v1.jinrishici.com/rensheng.txt"

	resp, err := http.Get(_url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return string(body), nil
}

// PraiseText 获取彩虹屁接口内容
func PraiseText() (text string, err error) {
	resp, err := http.Get("https://api.muxiaoguo.cn/api/caihongpi")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// 彩虹屁接口 JSON 结构体
	c := struct {
		Code interface{} `json:"code"`
		Msg  string      `json:"msg"`
		Data struct {
			Comment string `json:"comment"`
		} `json:"data"`
	}{}

	if err = json.Unmarshal(body, &c); err != nil {
		return
	}

	if v, ok := c.Code.(int); ok && v == -1 {
		return c.Msg, nil // request failed
	}

	return c.Data.Comment, nil
}
