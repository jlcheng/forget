package main

import (
	"github.com/jlcheng/forget/db"
	"github.com/rivo/tview"
	"fmt"
)

type TviewUI struct {
	app     *tview.Application
}

func main() {
	foo := db.Atlas{}
	bar := TviewUI{}
	fmt.Println(foo, bar)
}
