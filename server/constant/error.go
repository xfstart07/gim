// Author: xufei
// Date: 2019-09-27 15:23

package constant

import "github.com/pkg/errors"

var (
	ErrServerFail        = errors.New("server fail")
	ErrAccountRegistered = errors.New("用户已注册")
)
