package testkit

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
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