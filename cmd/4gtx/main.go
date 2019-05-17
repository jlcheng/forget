package main

import (
	"github.com/jlcheng/forget/cmd/4gtx/subcmd"
)
import _ "net/http/pprof"

func main() {
	subcmd.InitCobra()
	subcmd.Execute()
}
