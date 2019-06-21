package app

import (
	"encoding/json"
	"fmt"
	"github.com/jlcheng/forget/atlasrpc"
	"github.com/jlcheng/forget/cli"
	"strings"
)

// SnippetClient queries an Atlas server and renders the results using snippets
func SnippetClient(args []string) error {
	qterms := make([]string, len(args))
	for idx := range args {
		qterms[idx] = "+Body:" + args[idx]
	}
	sr, err := atlasrpc.RequestForBleveSearchResult(cli.Host(), cli.Port(), strings.Join(qterms, " "))
	if err != nil {
		return err
	}

	bytearr, err := json.Marshal(sr)
	if err != nil {
		return err
	}

	fmt.Println(string(bytearr))

	return nil
}
