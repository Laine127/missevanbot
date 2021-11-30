package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var botConfig *BotConfig

type BotConfig struct {
	Name   string        `mapstructure:"name"`   // 机器人的昵称
	Cookie string        `mapstructure:"cookie"` // 文件的存储位置
	Log    *LogConfig    `mapstructure:"log"`
	Redis  *RedisConfig  `mapstructure:"redis"` // Redis服务配置
	Push   *PushConfig   `mapstructure:"push"`  // 消息推送配置
	Admin  int           `mapstructure:"admin"` // 机器人控制人
	Rooms  []*RoomConfig `mapstructure:"rooms"` // 启用的房间列表配置
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	File       string `mapstructure:"file"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
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
	ID                 int    `mapstructure:"id"`                   // 直播间ID
	Name               string `mapstructure:"name"`                 // 主播自定义昵称
	Pinyin             bool   `mapstructure:"pinyin"`               // 是否开启注音功能
	RainbowMaxInterval int    `mapstructure:"rainbow_max_interval"` // 彩虹屁最长发送间隔
}

// LoadConfig is used to load configuration file
func LoadConfig() {
	conf := new(BotConfig)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error read configration file: %s \n", err))
	}
	if err := viper.Unmarshal(conf); err != nil {
		panic(fmt.Errorf("Fatal error unmarshal configration file: %s \n", err))
	}
	botConfig = conf
}

func Config() BotConfig {
	return *botConfig
}

func Cookie() string {
	return botConfig.Cookie
}
