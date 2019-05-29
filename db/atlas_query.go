package db

import (
	"bytes"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

// A query to Atlas returns an AtlasResponse, which is a collection of ResultEntry objects. AtlasResponse has the rough
// shape of:
//
//   AtlasResponse:
//   * ResultEntries // []ResultEntry
//     - NoteID
//     - Addr        // uint; Index to the first character of the line relative to the beginning of the Note.
//     - Line
//     * Spans       // []Span
//       - Start
//       - End
//
// (*: one-to-many; -: one-to-one)

type Span struct {
	Start uint
	End   uint
}

type ResultEntry struct {
	NoteID string
	Addr   uint
	Line   string
	Spans  []Span
}

type AtlasResponse struct {
	ResultEntries []ResultEntry
}

func (s *Atlas) QueryForBleveSearchResult(qstr string) (*bleve.SearchResult, error) {
	q := query.NewQueryStringQuery(qstr)
	sr := bleve.NewSearchRequest(q)
	sr.Fields = []string{"*"}
	sr.IncludeLocations = true
	return s.index.Search(sr)
}

/* == START: PrettyPrinter == */
func (atlasResponse *AtlasResponse) PPResultEntrySlice() string {
	var buf bytes.Buffer
	buf.WriteString("[]ResultEntry:\n")
	for idx, entry := range atlasResponse.ResultEntries {
		buf.WriteString(fmt.Sprintf("  entry[%d]:\n", idx))
		buf.WriteString(fmt.Sprintf("    NoteID: %s\n", entry.NoteID))
		buf.WriteString(fmt.Sprintf("    Addr: %d\n", entry.Addr))
		buf.WriteString(fmt.Sprintf("    Line: %s\n", entry.Line))
		buf.WriteString(fmt.Sprintf("    Spans: %v\n", entry.Spans))
	}
	buf.WriteString("\n")
	return buf.String()
}

/* == END: PrettyPrinter == */
