package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Name   string        `mapstructure:"name"`   // 机器人的昵称
	Cookie string        `mapstructure:"cookie"` // 文件的存储位置
	Redis  *RedisConfig  `mapstructure:"redis"`  // Redis服务配置
	Push   *PushConfig   `mapstructure:"push"`   // 消息推送配置
	Rooms  []*RoomConfig `mapstructure:"rooms"`  // 启用的房间列表配置
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"passwd"`
	DB       int    `mapstructure:"db"`
}

type PushConfig struct {
	Bark string `mapstructure:"bark"`
}

type RoomConfig struct {
	ID                 int      `mapstructure:"id"`                   // 直播间ID
	Name               string   `mapstructure:"name"`                 // 主播自定义昵称
	Pinyin             bool     `mapstructure:"pinyin"`               // 是否开启注音功能
	RainbowMaxInterval int      `mapstructure:"rainbow_max_interval"` // 彩虹屁最长发送间隔
	Rainbow            []string `mapstructure:"rainbow"`              // 彩虹屁自定义列表
}

// ReadConfig is used to read configuration file
func ReadConfig() {
	var conf Config
	// Read config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// 读取
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error read configration file: %s \n", err))
	}
	// 反序列化
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Fatal error unmarshal configration file: %s \n", err))
	}
	Conf = &conf
}
