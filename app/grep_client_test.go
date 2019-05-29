package app

import (
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/testkit"
	"testing"
	"reflect"
	"fmt"
)

const TEST_NOTE_1 = "the brown dog jumped over the red fox\na brown bird flew over the red fox\nthe brown chicken played\nwith the red hen\n...\nthe end"
const TEST_NOTE_2 = "the brown dog jumped over the red fox\na brown bird flew over brown-red fox\nthe brown chicken played\nwith the red hen\n...\nthe end"

func TestMapDocumentMatchToResultEntrySlice(t *testing.T) {
	testkit.DeleteTempIndexDir(t)
	tmpDir := testkit.GetTempIndexDir()
	atlas, err := db.Open(tmpDir, 2)
	if err != nil {
		t.Fatal(err)
	}

	note := db.Note{
		ID:         "test_note_1",
		Body:       TEST_NOTE_2,
		Title:      "",
		AccessTime: 0,
	}

	err = atlas.Enqueue(note)
	if err != nil {
		t.Fatal(err)
	}

	err = atlas.Flush()
	if err != nil {
		t.Fatal(err)
	}
	results, err := atlas.QueryForBleveSearchResult("brown")
	if err != nil {
		t.Fatal(err)
	}

	if len(results.Hits) == 0 {
		t.Fatal("search result empty")
	}
	dm := results.Hits[0]

	expected := []db.ResultEntry{
		{NoteID: "test_note_1", Line: "the brown dog jumped over the red fox", Addr: 0, Spans: []db.Span{{4, 9}}},
		{NoteID: "test_note_1", Line: "a brown bird flew over brown-red fox", Addr: 38, Spans: []db.Span{{2, 7}, {23, 28}}},
		{NoteID: "test_note_1", Line: "the brown chicken played", Addr: 75, Spans: []db.Span{{4, 9}}},
	}
	got, err := mapDocumentMatchToResultEntrySlice("Body", dm)
	if err != nil {
		t.Fatal("unexpected error.", err)
	}
	if !reflect.DeepEqual(expected, got) {
		fmt.Println("Unexpected result entries")
		fmt.Println("Expected:")
		fmt.Println(expected)
		fmt.Println("")
		fmt.Println("Got:")
		fmt.Println(got)
	}
}
