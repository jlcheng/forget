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
	"github.com/jlcheng/forget/watcher"
	"github.com/spf13/cobra"
	"log"
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
		fmt.Println("exwatch called with", exwatchArg.dataDir)

		// Creates a radovskyb.Watcher. Starts listening to its events. Finally, start the Watcherr.
		radwatcher := rwatch.New()
		if err := radwatcher.AddRecursive(exwatchArg.dataDir); err != nil {
			log.Fatalf("cannot watch %v: %v", exwatchArg, err)
		}
		stopCh := make(chan struct{})
		go watcher.ReceiveWatchEvents(radwatcher, stopCh)
		if err := radwatcher.Start(time.Millisecond * 100); err != nil {
			log.Fatal("cannot start watcher", err)
		}
		<-stopCh
	},
}

func init() {
	rootCmd.AddCommand(exwatchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exwatchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exwatchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	exwatchCmd.PersistentFlags().StringVar(&exwatchArg.dataDir, "dataDir", "", "data directory")
}
