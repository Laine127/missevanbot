package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Name   string  `mapstructure:"name"`   // 机器人的昵称
	Cookie string  `mapstructure:"cookie"` // 文件的存储位置
	Redis  *Redis  `mapstructure:"redis"`  // Redis服务配置
	Push   *Push   `mapstructure:"push"`   // 消息推送配置
	Rooms  []*Room `mapstructure:"rooms"`  // 启用的房间列表配置
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"passwd"`
	DB       int    `mapstructure:"db"`
}

type Push struct {
	Bark string `mapstructure:"bark"`
}

type Room struct {
	ID                 int      `mapstructure:"id"`                   // 直播间ID
	Name               string   `mapstructure:"name"`                 // 主播自定义昵称
	Pinyin             bool     `mapstructure:"pinyin"`               // 是否开启注音功能
	RainbowMaxInterval int      `mapstructure:"rainbow_max_interval"` // 彩虹屁最长发送间隔
	Rainbow            []string `mapstructure:"rainbow"`              // 彩虹屁自定义列表
}

// Load is used to load configuration file
func Load() {
	conf := new(Config)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error read configration file: %s \n", err))
	}
	if err := viper.Unmarshal(conf); err != nil {
		panic(fmt.Errorf("Fatal error unmarshal configration file: %s \n", err))
	}
	Conf = conf
}
