package subcmd

import (
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/trace"
	"github.com/jlcheng/forget/watcher"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var svrCmd = &cobra.Command{
	Use:   "svr",
	Short: "Runs 4gt server",
	Long:  `Runs 4gt server`,
	Run: func(cmd *cobra.Command, args []string) {
		runexsvr := watcher.NewWatcherFacade()
		defer runexsvr.Close()
		err := runexsvr.Listen(cli.Port(), cli.IndexDir(), cli.DataDirs(), time.Second*time.Duration(svrDuration))
		if err != nil {
			trace.PrintStackTrace(err)
			os.Exit(1)
		}
	},
}

var svrDuration int

func InitSvr() {
	rootCmd.AddCommand(svrCmd)
	svrCmd.Flags().IntVarP(&svrDuration, "duration", "t", 10, "seconds between polling fs for changes")
}
