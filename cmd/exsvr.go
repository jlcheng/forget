package cmd

import (
	"fmt"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/rpc"
	"github.com/jlcheng/forget/trace"
	"github.com/jlcheng/forget/watcher"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"

	rwatch "github.com/radovskyb/watcher"
)

var exsvrCmd = &cobra.Command{
	Use:   "exsvr",
	Short: "Server",
	Long: `Runs 4gt server`,
	Run: func(cmd *cobra.Command, args []string) {
		CliCfg.SetTraceLevel()
		closeCh := make(chan struct{})
		err := RunExsvr(exsvrPort, CliCfg.GetIndexDir(), CliCfg.DataDirs, closeCh)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var exsvrPort int
func init() {
	rootCmd.AddCommand(exsvrCmd)
	exsvrCmd.PersistentFlags().IntVarP(&exsvrPort, "port", "p", 8181, "rpc port")
}

func RunExsvr(port int, indexDir string, dataDirs []string, closeCh chan struct{}) error {
	atlas, err := db.Open(indexDir, 1)
	if err != nil {
		return err
	}
	docCount, err := atlas.GetDocCount()
	trace.Debug("atlas loc and size:", indexDir, docCount)

	fmt.Printf("Starting rpc on port %d\n", port)
	go rpc.StartRpcServer(atlas, port)

	// Creates a radovskyb.Watcher. Starts listening to its events. Finally, start the Watcher.
	radwatcher := rwatch.New()
	for _, dataDir := range dataDirs {
		_, err := os.Stat(dataDir)
		if err != nil {
			log.Fatalf("cannot watch %v: %v:\n", dataDir, err)
		}

		if err := radwatcher.AddRecursive(dataDir); err != nil {
			log.Fatalf("cannot watch %v: %v\n", dataDir, err)
		}
	}

	go watcher.ReceiveWatchEvents(atlas, radwatcher)

	if err := radwatcher.Start(time.Millisecond * 100); err != nil {
		log.Fatal("cannot start watcher", err)
	}

	// blocks until we receive a shutdown message
	select {
	case <-closeCh:
		radwatcher.Close()
	}

	return nil

}