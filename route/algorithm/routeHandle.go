// Author: xufei
// Date: 2019-11-27 11:33

package algorithm

type RouteHandle interface {
	RouteServer(values []string, key string) string
}
