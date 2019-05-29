package app

import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/atlasrpc"
	"github.com/jlcheng/forget/txtio"
	"strings"
	"time"
)

// GrepClient queries an Atlas server and renders the results similar to grep's output
func GrepClient(args []string) error {
	qterms := make([]string, len(args))
	for idx := range args {
		qterms[idx] = "+" + args[idx]
	}
	stime := time.Now()
	atlasResponse, err := atlasrpc.Request(cli.Host(), cli.Port(), strings.Join(args, " "))
	if err != nil {
		return err
	}
	fmt.Printf("Found %v notes in %v\n", len(atlasResponse.ResultEntries), time.Since(stime))
	for _, entry := range atlasResponse.ResultEntries {
		fmt.Println(txtio.AnsiFmt(entry))
	}
	return nil
}
