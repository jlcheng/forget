package db

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	"github.com/blevesearch/bleve/search/query"
	"github.com/jlcheng/forget/testkit"
	"testing"
)

const MOCK_NOTE_ONE = "the brown dog jumped over the red fox\na brown bird flew over the red fox\nthe brown chicken played\nwith the red hen\n...\nthe end"

func TestPrintDocumentMatch(t *testing.T) {
	testkit.DeleteTempIndexDir(t)
	tmpDir := testkit.GetTempIndexDir()
	atlas, err := Open(tmpDir, 2)
	if err != nil {
		t.Fatal(err)
	}
	note := Note{
		ID: "test_note_1",
		Body: MOCK_NOTE_ONE,
		Title: "",
		AccessTime:0,
	}
	atlas.Enqueue(note)
	atlas.Flush()
	index := atlas.rawIndex()
	q := query.NewQueryStringQuery("brown")
	sr := bleve.NewSearchRequest(q)
	sr.Fields = []string{"*"}
	sr.IncludeLocations = true
	results, err := index.Search(sr)
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Hits) == 0 {
		t.Fatal("search result empty")
	}
	termLocationMap := results.Hits[0].Locations["Body"]

	expected := "{\n  brown: [Location{Start:4, End:9, Pos:2},Location{Start:40, End:45, Pos:10},Location{Start:77, End:82, Pos:18}]\n}"
	if got := TermLocationToStr(&termLocationMap); got != expected {
		t.Fatal("unexpected formatting:", got)
	}
}

func MockNoteOneTermLocations() search.TermLocationMap {
	termLocationMap := search.TermLocationMap{}
	termLocationMap["brown"] = search.Locations{}
	locs := termLocationMap["brown"]
	locs = append(locs, &search.Location{Start:4, End:9, Pos:2})
	locs = append(locs, &search.Location{Start:40, End:45, Pos:10})
	locs = append(locs, &search.Location{Start:77, End:82, Pos:18})
	return termLocationMap
}

func TestDocumentMatchToPMatchInfo(t *testing.T) {
	dm := &search.DocumentMatch{}
	dm.Score = 12
	dm.Fields = make(map[string]interface{})
	dm.Fields["Body"] = MOCK_NOTE_ONE
	dm.Locations = make(map[string]search.TermLocationMap)
	dm.Locations["Body"] = MockNoteOneTermLocations()

	matchInfo := documentMatchToMatchInfo("Body", dm)
	_ = matchInfo
	fmt.Printf("%+v", dm)
}