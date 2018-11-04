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
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/trace"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

import _ "net/http/pprof"


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
		CliCfg.SetTraceLevel()

		if CliCfg.GetIndexDir() == "" {
			fmt.Println("index must be specified")
			return
		}
		if exloadArg.dataDir == "" {
			fmt.Println("dataDir must be specified")
			return
		}
		go func() {
			log.Println(http.ListenAndServe("localhost:8080", nil))
		}()

		err := CreateAndPopulateIndex(exloadArg.dataDir, CliCfg.GetIndexDir(), exloadArg.force)
		if err != nil {
			trace.OnError(err)
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

func CreateAndPopulateIndex(dataDir, indexDir string, force bool) error {
	trace.Debug(fmt.Sprintf("createAndPopulateIndex from (%v) to (%v)", dataDir, indexDir))

	// If indexDir exists, delete it or return error
	if f, err := os.Stat(indexDir); err == nil {
		if !force {
			return errors.New(fmt.Sprint("directory already exists:", indexDir))
		}
		if !f.IsDir() {
			return errors.New(fmt.Sprint("is a file:", indexDir))
		}
		trace.Debug("forcibly deleting indexDir:", indexDir)
		if err = os.RemoveAll(indexDir); err != nil {
			return err
		}
	}

	atlas, err := db.Open(indexDir, db.DEFAULT_BATCH_SIZE) // batch size of 1000
	if err != nil {
		return err
	}
	defer atlas.Close()

	stime := time.Now()
	helper := &indexHelper{atlas:atlas}

	dataDirInfo, err := os.Stat(dataDir)
	if err != nil {
		return err
	}

	if err = helper.indexFiles(dataDir, dataDirInfo); err != nil {
		return err
	}

	fmt.Printf("index stats: count: %v, total bytes: %v kb, elapsed time: %v\n",
		helper.count, helper.totalSize/1024, time.Since(stime))
	return nil
}

type indexHelper struct {
	atlas *db.Atlas
	totalSize int64
	count int
}

func (i *indexHelper) indexFiles(path string, info os.FileInfo) error {
	// Do not index .git
	if info.IsDir() && strings.HasSuffix(path, ".git") {
		trace.Debug(fmt.Sprintf("skipping %v", path))
		return filepath.SkipDir
	}

	// recurse into directory
	if info.IsDir() {
		cinfos, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}
		for _, childInfo := range cinfos {
			childPath := filepath.Join(path, childInfo.Name())
			cerr := i.indexFiles(childPath, childInfo)
			if cerr != nil && cerr != filepath.SkipDir {
				return cerr // bail on legitimate errors
			}
		}
	}

	// File name must contain a '.'
	if strings.LastIndexByte(path, '.') < strings.LastIndexByte(path, '/') {
		return nil
	}


	// Finally, index the heck out of this file
	doc := db.Note{
		ID: path,
		Body: debugReadFile(path),
		Title: info.Name(),
		AccessTime: info.ModTime().Unix(),
	}
	i.atlas.Enqueue(doc)
	i.totalSize = i.totalSize + info.Size()
	i.count = i.count + 1
	trace.Debug("indexed", doc.ID)
	return nil
}


func debugReadFile(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	defer f.Close()
	s, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return string(s)
}
