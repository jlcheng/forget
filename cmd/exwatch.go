// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/trace"
	"github.com/jlcheng/forget/watcher"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"

	rwatch "github.com/radovskyb/watcher"
)

var exwatchArg = struct{
	dataDir string
}{
	"",
}

// exwatchCmd represents the exwatch command
var exwatchCmd = &cobra.Command{
	Use:   "exwatch",
	Short: "Watches the target directory for changes",
	Long: `The exwatch command allows for protoyping of a file watcher feature`,
	Run: func(cmd *cobra.Command, args []string) {
		CliCfg.SetTraceLevel()

		dataDirs := viper.GetStringSlice(cli.DATA_DIRS)
		if len(dataDirs) == 0 {
			log.Fatal("dataDirs must be specified")
			return
		}
		fmt.Println("exwatch called with", dataDirs, CliCfg.GetIndexDir())

		atlas, err := db.Open(CliCfg.GetIndexDir(), BATCH_SIZE)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer atlas.Close()
		docCount, _:= atlas.GetDocCount()
		trace.Debug("atlas size:", docCount)

		// Creates a radovskyb.Watcher. Starts listening to its events. Finally, start the Watcher.
		radwatcher := rwatch.New()
		for _, dataDir := range CliCfg.GetDataDirs() {
			_, err := os.Stat(dataDir)
			if err != nil {
				log.Fatalf("cannot watch %v: %v\n:", dataDir, err)
			}

			if err := radwatcher.AddRecursive(dataDir); err != nil {
				log.Fatalf("cannot watch %v: %v\n", exwatchArg, err)
			}
		}

		wfacade := watcher.WatcherFacade{atlas}
		stopCh := make(chan struct{})
		go wfacade.ReceiveWatchEvents(radwatcher, stopCh)

		if err := radwatcher.Start(time.Millisecond * 100); err != nil {
			log.Fatal("cannot start watcher", err)
		}
		<-stopCh
	},
}

func init() {
	rootCmd.AddCommand(exwatchCmd)
}
