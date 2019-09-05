// Author: xufei
// Date: 2019-09-04 14:39

package server

import (
	"flag"

	"github.com/go-ini/ini"
)

type Config struct {
	ServerPort string `ini:"server_port"`
	RpcPort    string `ini:"rpc_port"`
	LogLevel   string `ini:"log_level"`
	Heartbeat  int    `json:"heartbeat"`
}

var (
	conf     *Config
	confPath string
)

func init() {
	flag.StringVar(&confPath, "config", "config/server.ini", "set server config filepath")
}

func GetConfig() *Config {
	return conf
}

func InitConfig() error {
	conf = defaultConfig()
	if err := readConfig(); err != nil {
		return err
	}
	return nil
}

func defaultConfig() *Config {
	return &Config{
		ServerPort: "8081",
		RpcPort:    "11211",
		LogLevel:   "info",
		Heartbeat:  30,
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
