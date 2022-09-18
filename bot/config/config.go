package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	BaseURL string // cqhttp 基础url
}

// Init 初始化配置文件
func (c *Config) Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	wd, _ := os.Getwd()
	viper.AddConfigPath(wd)
	if err := viper.ReadInConfig(); err != nil {
		panic("Read config file error")
	}
	c.BaseURL = viper.GetString("baseURL")
}
