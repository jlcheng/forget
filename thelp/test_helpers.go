package thelp

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"os"
	"path/filepath"
	"testing"
)

func TempIndex() (bleve.Index, error) {
	tmpDir := filepath.Join(os.TempDir(), "forget-test")
	err := os.RemoveAll(tmpDir)
	if err != nil {
		return nil, err
	}
	index, err := bleve.NewUsing(tmpDir, bleve.NewIndexMapping(), scorch.Name, scorch.Name, nil)
	if err != nil {
		return nil, err
	}
	return index, nil
}

func CleanUpIndex(t *testing.T, index bleve.Index) {
	if index == nil {
		return
	}
	err := index.Close()
	if err != nil {
		t.Error(err)
	}
}