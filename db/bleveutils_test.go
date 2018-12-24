package db

import (
	"github.com/blevesearch/bleve/search"
	"testing"
)

func TestTermLocationToStr(t *testing.T) {
	// brown: [Location{Start:4, End:9, Pos:2},Location{Start:40, End:45, Pos:10},Location{Start:77, End:82, Pos:18}]
	termLocationMap := search.TermLocationMap{}
	termLocationMap["brown"] = search.Locations{}
	locs := termLocationMap["brown"]
	locs = append(locs, &search.Location{Start:4, End:9, Pos:2})
	locs = append(locs, &search.Location{Start:40, End:45, Pos:10})
	locs = append(locs, &search.Location{Start:77, End:82, Pos:18})
	termLocationMap["brown"] = locs
	expected := "{\n  brown: [Location{Start:4, End:9, Pos:2},Location{Start:40, End:45, Pos:10},Location{Start:77, End:82, Pos:18}]\n}"
	if got := TermLocationToStr(&termLocationMap); got != expected {
		t.Fatal("unexpected formatting:", got)
	}
}