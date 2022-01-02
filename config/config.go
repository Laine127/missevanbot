package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	cfg GlobalConfig
)

type GlobalConfig struct {
	Admin int           `mapstructure:"admin"` // admin of the bot.
	Log   string        `mapstructure:"log"`   // log file path.
	Level string        `mapstructure:"level"` // log level.
	Redis *RedisConfig  `mapstructure:"redis"` // configurations of Redis client.
	Push  *PushConfig   `mapstructure:"push"`  // configurations of message push module.
	Rooms []*RoomConfig `mapstructure:"rooms"` // configurations of all the live rooms.
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
	ID      int    `mapstructure:"id"`     // ID of the room.
	Creator string `mapstructure:"name"`   // Nickname of the room creator.
	Enable  bool   `mapstructure:"enable"` // Whether to enable the bot.
	Watch   bool   `mapstructure:"watch"`  // Whether to enable live room status monitoring.
	Cookie  int    `mapstructure:"cookie"` // Cookie sequence number.
}

// Config return copy of the configurations.
func Config() GlobalConfig {
	return cfg
}

// Push return copy of the push configurations.
func Push() PushConfig {
	return *cfg.Push
}

// Admin return ID of the bot admin.
func Admin() int {
	return cfg.Admin
}

// LoadConfig is used to load configuration file
func LoadConfig() {
	cfg = GlobalConfig{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("data")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error read configration file: %s", err))
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("fatal error unmarshal configration file: %s", err))
	}
}
