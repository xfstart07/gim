// Author: xufei
// Date: 2019-09-04 14:39

package client

import (
	"flag"

	"github.com/go-ini/ini"
)

type Config struct {
	UserID         int64  `ini:"user_id"`
	Username       string `ini:"username"`
	WebPort        string `ini:"web_port"`
	ServerIP       string `ini:"server_ip"`
	ServerPort     string `ini:"server_port"`
	ServerRPCPort  string `ini:"server_rpc_port"`
	LogLevel       string `ini:"log_level"`
	ReconnectCount int    `ini:"reconnect_count"`
}

var (
	conf     *Config
	confPath string
)

func init() {
	flag.StringVar(&confPath, "config", "config/client.ini", "set client config filepath")
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
		UserID:         1434348343,
		Username:       "Leon",
		WebPort:        "8082",
		ServerIP:       "localhost",
		ServerPort:     "8083",
		ServerRPCPort:  "11211",
		LogLevel:       "info",
		ReconnectCount: 3,
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
