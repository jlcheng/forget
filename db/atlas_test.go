package db

import (
	"github.com/jlcheng/forget/testkit"
	"testing"
)

func TestAtlasOpen(t *testing.T) {
	testkit.DeleteTempIndexDir(t)
	tmpDir := testkit.GetTempIndexDir()
	atlas, err := Open(tmpDir, 2)
	if err != nil {
		t.Fatal(err)
	}

	got, err := atlas.GetDocCount()
	if err != nil {
		t.Fatal(err)
	}
	if got != 0 {
		t.Fatal("invalid docCount", got)
	}
}
