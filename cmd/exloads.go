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

	"github.com/spf13/cobra"
)

// exloadsCmd represents the exloads command
var exloadsCmd = &cobra.Command{
	Use:   "exloads",
	Short: "Loads files in a directory.",
	Long: `Loads files in a directory for prototyping

The exloads command will create a new index and populate it from a specified
directory. The document id will be the file name. The contents of the document
will be the body of the document (assumes UTF-8 encoding). The mtime will be
the timestamp of the document.

exloads will fail if the index path is non-empty.

exloads will not recurse inside the data directory - only the top level is used.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("exloads called")
	},
}

func init() {
	rootCmd.AddCommand(exloadsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exloadsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exloadsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
