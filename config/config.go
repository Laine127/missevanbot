package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

var (
	botConfig BotConfig
	cookie    string
)

type BotConfig struct {
	Cookie string        `mapstructure:"cookie"` // path of the cookie file.
	Admin  int           `mapstructure:"admin"`  // admin of the bot.
	Level  string        `mapstructure:"level"`  // log level.
	Redis  *RedisConfig  `mapstructure:"redis"`  // configurations of Redis client.
	Push   *PushConfig   `mapstructure:"push"`   // configurations of message push module.
	Rooms  []*RoomConfig `mapstructure:"rooms"`  // configurations of all the live rooms.
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
	ID                 int    `mapstructure:"id"`                   // ID of the room.
	Name               string `mapstructure:"name"`                 // nickname of the room creator.
	Enable             bool   `mapstructure:"enable"`               // whether to enable the bot.
	Pinyin             bool   `mapstructure:"pinyin"`               // whether to enable the pinyin function.
	RainbowMaxInterval int    `mapstructure:"rainbow_max_interval"` // the max interval for bait mode.
	Watch              bool   `mapstructure:"watch"`                // whether to enable live room status monitoring.
}

// Config return copy of the configurations.
func Config() BotConfig {
	return botConfig
}

// Push return copy of the push configurations.
func Push() PushConfig {
	return *botConfig.Push
}

// Cookie return the cookie in string format.
func Cookie() string {
	return cookie
}

// Admin return ID of the bot admin.
func Admin() int {
	return botConfig.Admin
}

// LoadConfig is used to load configuration file
func LoadConfig() {
	conf := BotConfig{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error read configration file: %s", err))
	}
	if err := viper.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("fatal error unmarshal configration file: %s", err))
	}
	botConfig = conf

	initCookie()
}

// initCookie read the content of cookie file which specified by botConfig.Cookie
// and store it in a global variable.
func initCookie() {
	file, err := os.Open(botConfig.Cookie)
	if err != nil {
		panic(fmt.Errorf("read cookie file failed: %s", err))
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("read cookie file failed: %s", err))
	}
	cookie = string(content)
}
