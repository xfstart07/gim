// Author: xufei
// Date: 2019-09-30 10:09

package ciface

type UserClient interface {
	Start() error
	Login()
	Shutdown()
}
