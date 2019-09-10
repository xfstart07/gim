// Author: xufei
// Date: 2019-09-09 16:16

package model

type ClientConfig struct {
	UserID         int64  `ini:"user_id"`
	Username       string `ini:"username"`
	WebPort        string `ini:"web_port"`
	ServerIP       string `ini:"server_ip"`
	ServerPort     string `ini:"server_port"`
	ServerRPCPort  string `ini:"server_rpc_port"`
	LogLevel       string `ini:"log_level"`
	ReconnectCount int    `ini:"reconnect_count"`
	HeartbeatTime  int    `ini:"heartbeat_time"`
	MsgLogPath     string `ini:"msg_log_path"`
}

type ServerConfig struct {
	ServerPort string `ini:"server_port"`
	RpcPort    string `ini:"rpc_port"`
	LogLevel   string `ini:"log_level"`
	Heartbeat  int    `ini:"heartbeat"`
	RedisURL   string `ini:"redis_url"`
	RedisPass  string `ini:"redis_pass"`
	RedisDB    int    `ini:"redis_db"`
}
