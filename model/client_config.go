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
}
