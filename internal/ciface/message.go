// Author: xufei
// Date: 2019-09-27 11:19

package ciface

type MessageHandler interface {
	InnerCommand(string) bool
	CheckMsg(string) bool
	SendMsg(string)
}
