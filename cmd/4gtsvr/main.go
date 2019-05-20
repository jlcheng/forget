package main

import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/trace"
	"github.com/jlcheng/forget/watcher"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"time"
)

import _ "net/http/pprof"

var duration = 10

var rootCmd = &cobra.Command{
	Use:   "4gt",
	Short: "Starts the 4gt server",
	Long:  `Starts the 4gt server`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.SetTraceLevel()

		// turn on pprof if specified
		if cli.PprofEnabled() {
			go func() {
				fmt.Println("pprof running on port 6060")
				log.Println(http.ListenAndServe("localhost:6060", nil))
			}()
		}

		runexsvr := watcher.NewWatcherFacade()
		defer runexsvr.Close()
		err := runexsvr.Listen(cli.Port(), cli.IndexDir(), cli.DataDirs(), time.Second*time.Duration(duration))
		if err != nil {
			trace.PrintStackTrace(err)
			os.Exit(1)
		}
	},
}

func main() {
	cobra.OnInitialize(cobraInit)

	cli.ConfigureFlagSet(rootCmd.Flags())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Loads configuration files
func cobraInit() {
	if err := cli.ProcessParsedFlagSet(rootCmd.Flags()); err != nil {
		log.Fatal(err)
	}
}
