package app

import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/rpc"
	"strings"
)

// SnippetClient queries an Atlas server and renders the results using snippets
func SnippetClient(args []string) error {
	qterms := make([]string, len(args))
	for idx := range args {
		qterms[idx] = "+" + args[idx]
	}
	sr, err := rpc.RequestForBleveSearchResult(cli.Host(), cli.Port(), strings.Join(args, " "))
	if err != nil {
		return err
	}
	fmt.Println(sr)
	return nil
}
