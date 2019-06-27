package subcmd

import (
	"github.com/jlcheng/forget/app"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/trace"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var svrCmd = &cobra.Command{
	Use:   "svr",
	Short: "Runs 4gt server",
	Long:  `Runs 4gt server`,
	Run: func(cmd *cobra.Command, args []string) {
		pollInterval := time.Second * time.Duration(svrPollInterval)
		err := app.StartAtlasServer(cli.Port(), cli.IndexDir(), cli.DataDirs(), svrRebuild, pollInterval)
		if err != nil {
			trace.PrintStackTrace(err)
			os.Exit(1)
		}
	},
}

var svrPollInterval int
var svrRebuild bool

func InitSvr() {
	rootCmd.AddCommand(svrCmd)
	svrCmd.Flags().IntVarP(&svrPollInterval, "duration", "t", 10, "seconds between polling fs for changes")
	svrCmd.Flags().BoolVar(&svrRebuild, "rebuild", false, "delete and rebuilds the index")
}
