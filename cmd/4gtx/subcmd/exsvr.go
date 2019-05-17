package subcmd

import (
	"fmt"
	"github.com/jlcheng/forget/trace"
	"github.com/jlcheng/forget/watcher"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var exsvrCmd = &cobra.Command{
	Use:   "exsvr",
	Short: "Server",
	Long:  `Runs 4gt server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command is deprecated. Use 4gtsvr instead")
		CliCfg.SetTraceLevel()

		runexsvr := watcher.NewWatcherFacade()
		defer runexsvr.Close()
		err := runexsvr.Listen(exsvrPort, CliCfg.GetIndexDir(), CliCfg.GetDataDirs(), time.Second*time.Duration(exsvrDuration))
		if err != nil {
			trace.PrintStackTrace(err)
			os.Exit(1)
		}
	},
}

var exsvrPort int
var exsvrDuration int

func InitExsvr() {
	rootCmd.AddCommand(exsvrCmd)
	exsvrCmd.PersistentFlags().IntVarP(&exsvrPort, "port", "p", 8181, "rpc port")
	exsvrCmd.PersistentFlags().IntVarP(&exsvrDuration, "duration", "t", 10, "seconds between polling fs for changes")
}
