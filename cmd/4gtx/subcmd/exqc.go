package subcmd

import (
	"fmt"
	"github.com/jlcheng/forget/rpc"
	"github.com/jlcheng/forget/txtio"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

var exqcCmd = &cobra.Command{
	Use:   "exqc",
	Short: "Client",
	Long:  `Runs 4gt exq network client`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command is deprecated. Use `4gt` instead")
		CliCfg.SetTraceLevel()

		qterms := make([]string, len(args))
		for idx := range args {
			qterms[idx] = "+" + args[idx]
		}
		stime := time.Now()
		atlasResponse, err := rpc.Request(exqcHost, exqcPort, strings.Join(args, " "))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found %v notes in %v\n", len(atlasResponse.ResultEntries), time.Since(stime))
		for _, entry := range atlasResponse.ResultEntries {
			fmt.Println(txtio.AnsiFmt(entry))
		}
		return

	},
}

var exqcPort int
var exqcHost string

func InitExqc() {
	rootCmd.AddCommand(exqcCmd)
	exqcCmd.PersistentFlags().IntVarP(&exqcPort, "port", "p", 8181, "rpc port")
	exqcCmd.PersistentFlags().StringVarP(&exqcHost, "host", "H", "localhost", "rpc host")
}
