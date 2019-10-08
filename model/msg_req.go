// Author: xufei
// Date: 2019-09-06 14:27

package model

type MsgReq struct {
	UserID int64  `json:"user_id"`
	Msg    string `json:"msg,omitempty"`
}
