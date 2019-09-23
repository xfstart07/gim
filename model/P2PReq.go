// Author: xufei
// Date: 2019-09-09 11:45

package model

import "strconv"

type P2PReq struct {
	UserID     int64  `json:"user_id"`
	ReceiverID int64  `json:"receiver_id"`
	Msg        string `json:"msg"`
}

func (p *P2PReq) ReceiverIDToString() string {
	return strconv.FormatInt(p.ReceiverID, 10)
}
