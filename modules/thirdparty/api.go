package thirdparty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// WeatherText 获取天气
func WeatherText(city string) (s string) {
	resp, err := http.Get(fmt.Sprintf("https://api.muxiaoguo.cn/api/tianqi?city=%s&type=1", city))
	if err != nil {
		log.Println("请求天气接口失败", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取天气接口失败", err.Error())
		return
	}
	// 天气接口 JSON 结构体
	c := new(struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			CityName string `json:"cityname"` // 城市名
			Weather  string `json:"weather"`  // 天气
			Temp     string `json:"temp"`     // 温度
			SD       string `json:"SD"`       // 相对湿度
		} `json:"data"`
	})
	if err := json.Unmarshal(body, c); err != nil {
		log.Println("解析天气接口失败", err.Error())
	}
	if c.Code != "200" || c.Msg != "success" {
		log.Println("请求天气接口错误", err.Error())
		return
	}
	return fmt.Sprintf(`%s 今日天气：
[天气] %s
[温度] %s
[相对湿度] %s`,
		c.Data.CityName,
		c.Data.Weather,
		c.Data.Temp,
		c.Data.SD,
	)
}

// PoemText 获取一句诗词
func PoemText() string {
	_url := "https://v1.jinrishici.com/rensheng.txt"

	resp, err := http.Get(_url)
	if err != nil {
		log.Println("获取诗词出错", err.Error())
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("获取诗词出错", err.Error())
		return ""
	}
	return string(body)
}

// PraiseText 获取彩虹屁接口内容
func PraiseText() (s string) {
	resp, err := http.Get("https://api.muxiaoguo.cn/api/caihongpi")
	if err != nil {
		log.Println("请求彩虹屁接口失败", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取彩虹屁接口失败", err.Error())
		return
	}
	// 彩虹屁接口 JSON 结构体
	c := new(struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Comment string `json:"comment"`
		} `json:"data"`
	})
	if err := json.Unmarshal(body, c); err != nil {
		log.Println("解析彩虹屁接口失败", err.Error())
	}
	if c.Code != "200" || c.Msg != "success" {
		log.Println("请求彩虹屁接口错误", c.Code)
		return
	}
	return c.Data.Comment
}
