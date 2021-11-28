package module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Weather 获取天气
func Weather(city string) (s string) {
	resp, err := http.Get(fmt.Sprintf("https://api.muxiaoguo.cn/api/tianqi?city=%s&type=1", city))
	if err != nil {
		ll.Print("请求天气接口失败", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ll.Print("读取天气接口失败", err.Error())
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
		ll.Print("解析天气接口失败", err.Error())
	}
	if c.Code != "200" || c.Msg != "success" {
		ll.Print("请求天气接口错误", c.Code)
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
