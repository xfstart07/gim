// Author: xufei
// Date: 2019-09-05 09:31

package main

import (
	"flag"
	"gim/server"
)

func main() {
	flag.Parse()

	srv := server.New()
	srv.Main()
}
