// Author: xufei
// Date: 2019-09-04 16:13

package model

import "fmt"

type User struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name" binding:"required"`
}

func (u *User) FormatMsg(msg string) string {
	return fmt.Sprintf("%s: 【%s】", u.UserName, msg)
}
