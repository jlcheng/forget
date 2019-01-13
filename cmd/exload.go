// Perform indexing on specified directories
package cmd

import (
	"errors"
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/trace"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

import _ "net/http/pprof"


var exloadArg = struct {
	force bool
}{
	false,
}

// exloadCmd represents the exloads command
var exloadCmd = &cobra.Command{
	Use:   "exload",
	Short: "Loads files in a directory.",
	Long: `Loads files in a directory for prototyping

The exload command will create a new index and populate it from specified directories. The document id will be the file 
name. The contents of the document will be the body of the document (assumes UTF-8 encoding). The mtime will be
the timestamp of the document.

exloads will fail if the index directory is non-empty.
`,
	Run: func(cmd *cobra.Command, args []string) {
		CliCfg.SetTraceLevel()
		dataDirs := viper.GetStringSlice(cli.DATA_DIRS)

		if CliCfg.GetIndexDir() == "" {
			fmt.Println("index must be specified")
			return
		}
		if len(dataDirs) == 0 {
			fmt.Println("dataDirs must be specified")
			return
		}
		err := CreateAndPopulateIndex(viper.GetStringSlice(cli.DATA_DIRS), CliCfg.GetIndexDir(), exloadArg.force)
		if err != nil {
			trace.OnError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(exloadCmd)
	exloadCmd.PersistentFlags().BoolVarP(&exloadArg.force, "force", "f", false, "forces index to run, even if indexDir already exists")
}

func CreateAndPopulateIndex(dataDirs []string, indexDir string, force bool) error {
	trace.Debug(fmt.Sprintf("createAndPopulateIndex from (%v) to (%v)",dataDirs, indexDir))

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

	for _, dataDir := range dataDirs {
		dataDirInfo, err := os.Stat(dataDir)
		if err != nil {
			return err
		}

		if err = helper.indexFiles(dataDir, dataDirInfo); err != nil {
			return err
		}
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
	if !db.FilterFile(path, info) {
		return nil
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
