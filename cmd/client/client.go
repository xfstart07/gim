// Author: xufei
// Date: 2019-09-04

package main

import (
	"flag"
	"gim/client"
)

func main() {
	flag.Parse()

	cli := client.New()
	cli.Main()
}
