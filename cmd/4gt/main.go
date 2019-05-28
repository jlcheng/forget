package main

import (
	"github.com/jlcheng/forget/cmd/4gt/subcmd"
)

func main() {
	subcmd.InitRoot()
	subcmd.Execute()
}
