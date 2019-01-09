package cmd

import (
	"fmt"
	"github.com/jlcheng/forget/rpc"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var exqcCmd = &cobra.Command{
	Use:   "exqc",
	Short: "Client",
	Long: `Runs 4gt exq network client`,
	Run: func(cmd *cobra.Command, args []string) {
		CliCfg.SetTraceLevel()

		qterms := make([]string, len(args))
		for idx, _ := range args {
			qterms[idx] = "+" + args[idx]
		}
		response, err := rpc.Request(exqcHost, exqcPort, strings.Join(args, " "))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(response)
	},
}

var exqcPort int
var exqcHost string
func init() {
	rootCmd.AddCommand(exqcCmd)
	exqcCmd.PersistentFlags().IntVarP(&exqcPort, "port", "p", 8181, "rpc port")
	exqcCmd.PersistentFlags().StringVarP(&exqcHost, "host", "H", "localhost", "rpc host")
}
