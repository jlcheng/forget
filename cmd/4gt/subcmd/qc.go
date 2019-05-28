package subcmd

import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/rpc"
	"github.com/jlcheng/forget/txtio"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

var qcCmd = &cobra.Command{
	Use:   "qc",
	Short: "Client",
	Long:  `Runs 4gt network query client`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.SetTraceLevel()

		qterms := make([]string, len(args))
		for idx := range args {
			qterms[idx] = "+" + args[idx]
		}
		stime := time.Now()
		atlasResponse, err := rpc.Request(cli.Host(), cli.Port(), strings.Join(args, " "))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found %v notes in %v\n", len(atlasResponse.ResultEntries), time.Since(stime))
		for _, entry := range atlasResponse.ResultEntries {
			fmt.Println(txtio.AnsiFmt(entry))
		}
	},
}

func InitExqc() {
	rootCmd.AddCommand(qcCmd)
}
