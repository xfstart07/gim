// Author: xufei
// Date: 2019-09-09 16:16

package model

type ClientConfig struct {
	User

	WebPort    string `ini:"web_port"`
	ServerIP   string `ini:"server_ip"`
	ServerPort string `ini:"server_port"`

	LogLevel       string `ini:"log_level"`
	MsgLogPath     string `ini:"msg_log_path"`
	ReconnectCount int    `ini:"reconnect_count"`
	HeartbeatTime  int    `ini:"heartbeat_time"`

	EtcdUrl        string `ini:"etcd_url"`
	EtcdServerName string `ini:"etcd_server_name"`
}

type ServerConfig struct {
	WebEnable  bool   `ini:"web_enable"`
	ServerPort string `ini:"server_port"`
	RpcPort    string `ini:"rpc_port"`
	Heartbeat  int    `ini:"heartbeat"`

	LogLevel string `ini:"log_level"`

	RedisURL  string `ini:"redis_url"`
	RedisPass string `ini:"redis_pass"`
	RedisDB   int    `ini:"redis_db"`

	EtcdUrl        string `ini:"etcd_url"`
	EtcdServerName string `ini:"etcd_server_name"`
}
