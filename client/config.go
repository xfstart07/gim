// Author: xufei
// Date: 2019-09-04 14:39

package client

import (
	"flag"
	"gim/internal/lg"
	"gim/model"

	"github.com/go-ini/ini"
)

var (
	conf         *model.ClientConfig
	confPath     string
	confUserID   int64
	confUserName string
)

func init() {
	flag.StringVar(&confPath, "config", "config/client.ini", "set client config filepath")
	flag.Int64Var(&confUserID, "user_id", 0, "userInfo id")
	flag.StringVar(&confUserName, "username", "", "userInfo name")
}

func GetConfig() *model.ClientConfig {
	return conf
}

func InitConfig() error {
	conf = defaultConfig()
	if err := readConfig(); err != nil {
		return err
	}

	lg.Logger().Sugar().Info(conf)

	// 根据命令行传入信息更新用户信息
	if confUserID != 0 {
		conf.UserID = confUserID
		conf.UserName = confUserName
	}

	return nil
}

func defaultConfig() *model.ClientConfig {
	return &model.ClientConfig{
		WebPort:        "8082",
		ServerURL:      "http://localhost:8081",
		LogLevel:       "info",
		ReconnectCount: 3,
		HeartbeatTime:  60,
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
