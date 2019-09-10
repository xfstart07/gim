// Author: xufei
// Date: 2019-09-10 17:11

package model

type PushMsg struct {
	UserID  int64  `json:"user_id"`
	Msg     string `json:"msg"`
	MsgType int32  `json:"msg_type"`
}
