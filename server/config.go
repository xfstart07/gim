// Author: xufei
// Date: 2019-09-04 14:39

package server

import (
	"flag"
	"gim/model"

	"github.com/go-ini/ini"
)

var (
	conf     *model.ServerConfig
	confPath string
)

func init() {
	flag.StringVar(&confPath, "config", "config/server.ini", "set server config filepath")
}

func GetConfig() *model.ServerConfig {
	return conf
}

func InitConfig() error {
	conf = defaultConfig()
	if err := readConfig(); err != nil {
		return err
	}
	return nil
}

func defaultConfig() *model.ServerConfig {
	return &model.ServerConfig{
		WebEnable:  true,
		ServerPort: "8081",
		RpcPort:    "11211",
		LogLevel:   "info",
		Heartbeat:  30,
		RedisURL:   "localhost:6379",
		RedisPass:  "",
		RedisDB:    0,
	}
}

func readConfig() error {
	file, err := ini.Load(confPath)
	if err != nil {
		return err
	}

	if err := file.MapTo(&conf); err != nil {
		return err
	}
	return nil
}
