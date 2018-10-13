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
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/trace"
	"github.com/spf13/cobra"
	"log"
)

// exdumpCmd represents the exdump command
var exdumpCmd = &cobra.Command{
	Use:   "exdump",
	Short: "Dumps the index into the specified directory.",
	Long:  "Dumps the index into the specified directory.",
	Run: func(cmd *cobra.Command, args []string) {
		setDebugLevel()
		fmt.Println("exdump called")
		atlas, err := db.Open(indexDir)
		if err != nil {
			log.Fatal(err)
		}
		if err = IterateDocuments(atlas, nil); err != nil {
			log.Fatal(err)
		}
		atlas.Close()
	},
}

func init() {
	rootCmd.AddCommand(exdumpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exdumpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exdumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type IDHandler func(param ...interface{}) error
// IterateDocuments
func IterateDocuments(atlas *db.Atlas, foo IDHandler) error {
	trace.Debug("IterateDocuments")
	docs, err := atlas.DumpAll()
	if err != nil {
		return err
	}
	for idx, doc := range docs {
		trace.Debug(doc.PrettyString())
		if foo != nil {
			foo(idx, doc)
		}
	}
	return nil
}

