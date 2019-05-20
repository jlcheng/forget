package main

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

// This must be global so we can bind values parsed by Cobra to the Viper API
var rootCmd = &cobra.Command{
	Use:   "4gt",
	Short: "Query a 4gt server",
	Long:  "Query a 4gt server",
	Run:   cobraMain,
}

func cobraMain(cmd *cobra.Command, args []string) {
	cli.SetTraceLevel()

	qterms := make([]string, len(args))
	for idx, elem := range args {
		if strings.Contains(elem, " ") {
			qterms[idx] = fmt.Sprintf("+Body:\"%v\"", elem)
		} else {
			qterms[idx] = fmt.Sprintf("+Body:%v", elem)
		}
	}
	SendQuery(cli.Host(), cli.Port(), strings.Join(qterms, " "))
}

// Loads configuration files
func cobraInit() {
	if err := cli.ProcessParsedFlagSet(rootCmd.Flags()); err != nil {
		log.Fatal(err)
	}
}

// TODO: Move this into the main directory
func SendQuery(host string, port int, query string) {
	stime := time.Now()
	atlasResponse, err := rpc.Request(host, port, query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %v lines in %v\n", len(atlasResponse.ResultEntries), time.Since(stime))
	for _, entry := range atlasResponse.ResultEntries {
		fmt.Println(txtio.AnsiFmt(entry))
	}
}

func main() {
	cobra.OnInitialize(cobraInit)

	cli.ConfigureFlagSet(rootCmd.Flags())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
