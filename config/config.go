package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Cookie string        `mapstructure:"cookie"`
	Redis  *RedisConfig  `mapstructure:"redis"`
	Push   *PushConfig   `mapstructure:"push"`
	Rooms  []*RoomConfig `mapstructure:"rooms"`
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
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	RainbowInterval int      `json:"rainbow_interval"`
	Rainbow         []string `json:"rainbow"`
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
