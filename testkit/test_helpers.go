package testkit

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/blevesearch/bleve/mapping"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	TEST_TMP_IDX_DIR = "forget-test-index"
)

func DeleteTempIndexDir(t *testing.T) {
	tmpDir := GetTempIndexDir()
	err := os.RemoveAll(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
}

func GetTempIndexDir() string {
	return filepath.Join(os.TempDir(), TEST_TMP_IDX_DIR)
}

func MkTempIndex(t *testing.T) bleve.Index {
	tmpDir := GetTempIndexDir()
	err := os.RemoveAll(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	index, err := bleve.NewUsing(tmpDir, bleve.NewIndexMapping(), scorch.Name, scorch.Name, nil)
	if err != nil {
		t.Fatal(err)
	}
	return index
}

func CleanUpTempIndex(t *testing.T, index bleve.Index) {
	if index == nil {
		return
	}
	err := index.Close()
	if err != nil {
		t.Fatal(err)
	}
	tmpDir := GetTempIndexDir()
	err = os.RemoveAll(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
}

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
	defer StopOnError(os.RemoveAll(tmpDir))
	defer StopOnError(index.Close())

	err = tiCtx(index, tmpDir)
	if err != nil {
		return err
	}
	return nil
}

func StopOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}