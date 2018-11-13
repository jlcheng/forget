package db

import (
	"github.com/jlcheng/forget/testkit"
	"reflect"
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

func TestAtlasRead(t *testing.T) {

	testCases := []struct{
		GivenNotes []Note
		Query string
		Notes []Note
	}{
		{
			GivenNotes: []Note{
				{ID: "ID", Body: "Body"},
			},
			Query: "Body",
			Notes: []Note{
				{ID:"ID", Body:"Body"},
			},
		},
		{
			GivenNotes: []Note{
				{ID: "ID", Body: "booty"},
			},
			Query: "Body",
			Notes: []Note{
			},
		},
		{
			GivenNotes: []Note{
				{ID: "ID", Body: "red fox jumps over the brown dog"},
			},
			Query: "brown",
			Notes: []Note{
				{ID: "ID", Body: "red fox jumps over the brown dog"},
			},
		},
	}

	for idx, tcase := range testCases {
		testkit.DeleteTempIndexDir(t)
		tmpDir := testkit.GetTempIndexDir()
		atlas, err := Open(tmpDir, 2)
		if err != nil {
			t.Fatal(err)
		}

		for _, given := range tcase.GivenNotes {
			atlas.Enqueue(given)
		}
		atlas.Flush()

		actualNotes, err := atlas.QueryString(tcase.Query)
		if err != nil {
			t.Fatal(err, "test case: ", idx)
		}
		if !reflect.DeepEqual(tcase.Notes, actualNotes) {
			t.Fatal("comparison failed. test case:", idx, ",", tcase.Notes, actualNotes)
		}
	}
}
