package testkit

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/blevesearch/bleve/mapping"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	TEST_TMP_IDX_DIR = "forget-test-index"
)

// DeleteTempIndexDir removes the temporary directory used in tests
func DeleteTempIndexDir(t *testing.T) {
	tmpDir := GetTempIndexDir()
	err := os.RemoveAll(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
}

// GetTempIndexDir returns the path to a temporary test directory. The directory is not safe for concurrent use.
func GetTempIndexDir() string {
	return filepath.Join(os.TempDir(), TEST_TMP_IDX_DIR)
}

// TempIndexContest is a function that acts on a Bleve index with tmpDir as a temporary directory
type TempIndexContext func(index bleve.Index, tmpDir string) error

func DoInTempIndexContext(tiCtx TempIndexContext, indexMapping mapping.IndexMapping) error {
	tmpDir := GetTempIndexDir()
	err := os.RemoveAll(tmpDir)
	if err != nil {
		return err
	}

	index, err := bleve.NewUsing(tmpDir, indexMapping, scorch.Name, scorch.Name, nil)
	if err != nil {
		return err
	}
	defer TempDirRemoveAll(tmpDir)
	defer TempCloserClose(index)

	err = tiCtx(index, tmpDir)
	if err != nil {
		return err
	}
	return nil
}

// LogError consumes an error by logging it.
func TempDirRemoveAll(path string) {
	if err := os.RemoveAll(path); err != nil {
		log.Fatal(err)
	}
}

func TempCloserClose(f io.Closer) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}