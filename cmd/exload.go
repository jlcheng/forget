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
	"errors"
	"fmt"
	"github.com/jlcheng/forget/log"
	"github.com/jlcheng/forget/globals"
	"github.com/jlcheng/forget/search"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var exloadArg = struct {
	dataDir string
	force bool
}{
	"",
	false,
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
		switch logLevelStr {
		case "DEBUG":
			log.Level = log.LOG_DEBUG
		case "WARN":
			log.Level = log.LOG_WARN
		default:
			log.Level = log.LOG_NONE
		}
		if indexDir == "" {
			log.Warn("index must be specified")
			return
		}
		if exloadArg.dataDir == "" {
			log.Warn("dataDir must be specified")
			return
		}
		log.Debug("exload called with args:", args)
		log.Debug("exload called with indexDir:", indexDir)
		log.Debug("exload called with dataDir:", exloadArg.dataDir)
		err := CreateAndPopulateIndex(exloadArg.dataDir, indexDir, exloadArg.force)
		if err != nil {
			log.OnError(err)
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
	exloadCmd.PersistentFlags().BoolVarP(&exloadArg.force, "force", "f", false, "forces index to run, even if indexDir already exists")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// rename application variables
var searchEngine = globals.SearchEngine

func CreateAndPopulateIndex(dataDir, indexDir string, force bool) error {
	log.Debug(fmt.Sprintf("createAndPopulateIndex from %v to %v", dataDir, indexDir))
	f, err := os.Stat(indexDir)
	if err == nil {
		if !force {
			return errors.New(fmt.Sprint("directory already exists:", indexDir))
		}
		// delete and recreate index
		if !f.IsDir() {
			return errors.New(fmt.Sprint("is a file:", indexDir))
		}
		log.Debug("forcibly deleting indexDir:", indexDir)
		if err = os.Remove(indexDir); err != nil {
			return err
		}
	}

	if err = os.Mkdir(indexDir, 0755); err != nil {
		return err
	}

	log.Debug("index directory created, starting to index")

	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		doc := search.Document{
			Id: file.Name(),
			Body: file.Name(),
			AccessTime: file.ModTime(),
		}

		searchEngine.Enqueue(doc)
	}

	return nil
}