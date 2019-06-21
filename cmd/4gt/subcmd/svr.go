package subcmd

import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/db/files"
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
		if rebuild {
			err := os.RemoveAll(cli.IndexDir())
			if err != nil {
				trace.Warn(fmt.Sprintf("cannot recursively delete %v", cli.IndexDir()))
				return
			}

			atlas, err := db.Open(cli.IndexDir(), 100)
			if err != nil {
				trace.Warn(err)
				return
			}
			err = files.RebuildIndex(atlas, cli.DataDirs())
			if err != nil {
				trace.Warn(err)
			}
			err = atlas.Flush()
			if err != nil {
				trace.Warn(err)
			}
			err = atlas.Close()
			if err != nil {
				trace.Warn(err)
			}
		}

		err := runexsvr.Listen(cli.Port(), cli.IndexDir(), cli.DataDirs(), time.Second*time.Duration(svrDuration))
		if err != nil {
			trace.PrintStackTrace(err)
			os.Exit(1)
		}
	},
}

var svrDuration int
var rebuild bool

func InitSvr() {
	rootCmd.AddCommand(svrCmd)
	svrCmd.Flags().IntVarP(&svrDuration, "duration", "t", 10, "seconds between polling fs for changes")
	svrCmd.Flags().BoolVar(&rebuild, "rebuild", false, "delete and rebuild the index")
}
