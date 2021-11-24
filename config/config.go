package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Cookie string     `mapstructure:"cookie"`
	Redis  *RedisConf `mapstructure:"redis"`
}

type RedisConf struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"passwd"`
	DB       int    `mapstructure:"db"`
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
