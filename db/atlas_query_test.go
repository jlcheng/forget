package db

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	"github.com/blevesearch/bleve/search/query"
	"github.com/jlcheng/forget/testkit"
	"reflect"
	"testing"
)

const TEST_NOTE_1 = "the brown dog jumped over the red fox\na brown bird flew over the red fox\nthe brown chicken played\nwith the red hen\n...\nthe end"
const TEST_NOTE_2 = "the brown dog jumped over the red fox\na brown bird flew over brown-red fox\nthe brown chicken played\nwith the red hen\n...\nthe end"

func TestPrintDocumentMatch(t *testing.T) {
	testkit.DeleteTempIndexDir(t)
	tmpDir := testkit.GetTempIndexDir()
	atlas, err := Open(tmpDir, 2)
	if err != nil {
		t.Fatal(err)
	}
	note := Note{
		ID:         "test_note_1",
		Body:       TEST_NOTE_1,
		Title:      "",
		AccessTime: 0,
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

func MockNoteOneTermLocationMap() search.TermLocationMap {
	termLocationMap := search.TermLocationMap{}
	termLocationMap["brown"] = search.Locations{}
	locs := termLocationMap["brown"]
	locs = append(locs, &search.Location{Start:4, End:9, Pos:2})
	locs = append(locs, &search.Location{Start:40, End:45, Pos:10})
	locs = append(locs, &search.Location{Start:77, End:82, Pos:18})
	termLocationMap["brown"] = locs
	return termLocationMap
}

func TestGetLineAround(t *testing.T) {
	var text, gotL, expectL string
	var gotLno, expectLno uint64

	// \n  .  .  b  o  d  y  0  .  . \n
	//  6  7  8  9 10 11 12 13 14 15 16
	text = "header\n..body0..\nfooter"
	expectLno, expectL = 7, "..body0.."
	gotLno, gotL = getLineAround(text, 9, 14)
	if gotLno != expectLno {
		t.Fatal("unexpected line number:", gotLno)
	}
	if gotL != expectL {
		t.Fatalf("unexpected line [%s]\n", gotL)
	}

	// \n  b  o  d  y  1  .  . \n
	//  6  7  8  9 10 11 12 13 14 15 16
	text = "header\nbody1..\nfooter"
	expectLno, expectL = 7, "body1.."
	gotLno, gotL = getLineAround(text, 7, 11)
	if gotLno != expectLno {
		t.Fatal("unexpected line number:", gotLno)
	}
	if gotL != expectL {
		t.Fatalf("unexpected line [%s]\n", gotL)
	}

	// \n  .  .  b  o  d  y  2 \n
	//  6  7  8  9 10 11 12 13 14 15 16
	text = "header\n..body2\nfooter"
	expectLno, expectL = 7, "..body2"
	gotLno, gotL = getLineAround(text, 9, 14)
	if gotLno != expectLno {
		t.Fatal("unexpected line number:", gotLno)
	}
	if gotL != expectL {
		t.Fatalf("unexpected line [%s]\n", gotL)
	}

	//  .  .  b  o  d  y  3  .  . \n
	//  0  1  2  3  4  5  6  7  8  9
	text = "..body3..\nfooter"
	expectLno, expectL = 0, "..body3.."
	gotLno, gotL = getLineAround(text, 2, 7)
	if gotLno != expectLno {
		t.Fatal("unexpected line number:", gotLno)
	}
	if gotL != expectL {
		t.Fatalf("unexpected line [%s]\n", gotL)
	}

	// \n  .  .  b  o  d  y  4  .  .
	//  6  7  8  9 10 11 12 13 14 15
	text = "header\n..body4.."
	expectLno, expectL = 7, "..body4.."
	gotLno, gotL = getLineAround(text, 9, 14)
	if gotLno != expectLno {
		t.Fatal("unexpected line number:", gotLno)
	}
	if gotL != expectL {
		t.Fatalf("unexpected line [%s]\n", gotL)
	}

	// \n  .  .  b  o  d  y  5
	//  6  7  8  9 10 11 12 13
	text = "header\n..body5"
	expectLno, expectL = 7, "..body5"
	gotLno, gotL = getLineAround(text, 9, 14)
	if gotLno != expectLno {
		t.Fatal("unexpected line number:", gotLno)
	}
	if gotL != expectL {
		t.Fatalf("unexpected line [%s]\n", gotL)
	}
}

func TestMapDocumentMatch(t *testing.T) {
	dm := &search.DocumentMatch{}
	dm.Score = 12
	dm.Fields = make(map[string]interface{})
	dm.Fields["Body"] = TEST_NOTE_1
	dm.Locations = make(map[string]search.TermLocationMap)
	dm.Locations["Body"] = MockNoteOneTermLocationMap()

	srLocations, ok := mapDocumentMatch("Body", dm)
	if !ok {
		t.Fatal("mapping failed")
	}

	if srLocations.NoteID != dm.ID {
		t.Fatal("unexpected NoteID")
	}
	if srLocations.Body != dm.Fields["Body"] {
		t.Fatal("unexpected Body")
	}
	expected := dm.Locations["Body"]
	got := srLocations.TermLocationMap
	if !reflect.DeepEqual(got, expected) {
		t.Fatal("unexpected Locations['Body']")
	}
}

func TestQueryForMatches(t *testing.T) {
	testkit.DeleteTempIndexDir(t)
	tmpDir := testkit.GetTempIndexDir()
	atlas, err := Open(tmpDir, 2)
	if err != nil {
		t.Fatal(err)
	}
	note := Note{
		ID:         "test_note_1",
		Body:       TEST_NOTE_1,
		Title:      "",
		AccessTime: 0,
	}
	atlas.Enqueue(note)
	atlas.Flush()
	queryMatches, err := atlas.QueryForMatches("fox")
	if err != nil {
		t.Fatal(err)
	}
	if len(queryMatches.MatchInfoMap) == 0 {
		t.Fatal("unexpected size")
	}
}

func TestQueryMatchNoteIDs(t *testing.T) {
	qm := QueryMatches{}
	qm.MatchInfoMap = make(map[string]*MatchInfo)
	qm.MatchInfoMap["0003"] = nil
	qm.MatchInfoMap["0001"] = nil
	qm.MatchInfoMap["0002"] = nil

	expected := []string{"0001", "0002", "0003"}
	if got := qm.NoteIDs(); !reflect.DeepEqual(expected, got) {
		t.Fatal("unexpected NoteIDs:", got)
	}
}

func TestMapDocumentMatchToResultEntrySlice(t *testing.T) {
	testkit.DeleteTempIndexDir(t)
	tmpDir := testkit.GetTempIndexDir()
	atlas, err := Open(tmpDir, 2)
	if err != nil {
		t.Fatal(err)
	}
	note := Note{
		ID:         "test_note_1",
		Body:       TEST_NOTE_2,
		Title:      "",
		AccessTime: 0,
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
	dm := results.Hits[0]

	expected := []ResultEntry {
		{NoteID: "test_note_1", Line: "the brown dog jumped over the red fox", Addr: 0, Spans:[]Span {{4, 9}}},
		{NoteID: "test_note_1", Line: "a brown bird flew over brown-red fox", Addr: 38, Spans:[]Span {{2, 7}, {23, 28}}},
		{NoteID: "test_note_1", Line: "the brown chicken played", Addr: 75, Spans:[]Span {{4, 9}}},
	}
	got, err := mapDocumentMatchToResultEntrySlice("Body", dm)
	if err != nil {
		t.Fatal("unexpected error.", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Fatal("unexpected result entries.", got)
	}
}