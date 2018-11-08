package testkit

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"os"
	"path/filepath"
	"testing"
)

func MkTempIndex(t *testing.T) bleve.Index {
	tmpDir := filepath.Join(os.TempDir(), "forget-test")
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
	tmpDir := filepath.Join(os.TempDir(), "forget-test")
	err = os.RemoveAll(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

}