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
	"log"
	"os"
	"errors"
	"io/ioutil"
)

var exloadArg = struct {
	dataDir string
}{
	"",
}

// exloadCmd represents the exloads command
var exloadCmd = &cobra.Command{
	Use:   "exload",
	Short: "Loads files in a directory.",
	Long: `Loads files in a directory for prototyping

The exload command will create a new index and populate it from a specified
directory. The document id will be the file name. The contents of the document
will be the body of the document (assumes UTF-8 encoding). The mtime will be
the timestamp of the document.

exloads will fail if the index path is non-empty.

exloads will not recurse inside the data directory - only the top level is used.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if indexDir == "" {
			fmt.Println("index must be specified")
			return
		}
		if exloadArg.dataDir == "" {
			fmt.Println("dataDir must be specified")
			return
		}
		fmt.Printf("exload called with args: %v\n", args)
		fmt.Printf("exload called with indexDir: %v\n", indexDir)
		fmt.Println("exload called with dataDir: ", exloadArg.dataDir)
		err := CreateAndPopulateIndex(exloadArg.dataDir, indexDir)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(exloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exloadCmd.PersistentFlags().String("foo", "", "A help for foo")
	exloadCmd.PersistentFlags().StringVar(&exloadArg.dataDir, "dataDir", "", "data directory")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func CreateAndPopulateIndex(dataDir, indexDir string) error {
	log.Printf("createAndPopulateIndex from %v to %v\n", dataDir, indexDir)
	_, err := os.Stat(indexDir)
	if err == nil {
		return errors.New(fmt.Sprint("directory already exists: ", indexDir))
	}

	err = os.Mkdir(indexDir, 0755)
	if err != nil {
		return err
	}

	log.Println("index directory created, starting to index")

	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		return err
	}
	for idx, file := range files {
		log.Println(idx, file.Name(), file.IsDir(), file.Size())
	}

	return nil
}