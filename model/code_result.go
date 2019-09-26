// Author: xufei
// Date: 2019-09-04 16:16

package model

const (
	CodeSuccessed = "0"
)

type CodeResult struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrResult struct {
	Message string     `json:"message"`
	Errors  []ErrField `json:"errors,omitempty"`
}

type ErrField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
