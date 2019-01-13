package db

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	"github.com/blevesearch/bleve/search/query"
	"log"
	"strings"
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
	End uint
}

type ResultEntry struct {
	NoteID string
	Addr uint
	Line string
	Spans []Span
}

type AtlasResponse struct {
	ResultEntries []ResultEntry
}

func (s *Atlas) QueryForResponse(qstr string) AtlasResponse {
	q := query.NewQueryStringQuery(qstr)
	sr := bleve.NewSearchRequest(q)
	sr.SortBy([]string{ACCESS_TIME})
	sr.Fields = []string{"*"}
	sr.IncludeLocations = true
	searchResult, err := s.index.Search(sr)
	if err != nil {
		return AtlasResponse{ResultEntries: make([]ResultEntry, 0)}
	}
	return mapSearchResult("Body", searchResult)
}

func getLineAround(text string, start, end uint64) (uint64, string) {
	lineStart := strings.LastIndexByte(text[:start], '\n') + 1
	lineEnd := strings.IndexRune(text[end:], '\n')
	if lineEnd != -1 {
		lineEnd = int(end) + lineEnd
	} else {
		lineEnd = len(text)
	}
	line := text[lineStart:lineEnd]
	return uint64(lineStart), line
}

// Flattens the DocumentMatch's Term=>Location map into a collection of matching lines. The error is non-nil if
// mapping fails.
func mapDocumentMatchToResultEntrySlice(fieldName string, dm *search.DocumentMatch) ([]ResultEntry, error) {
	emptyResponse := make([]ResultEntry, 0, 0)
	body, err := resolveBody(fieldName, dm)
	if err != nil {
		return emptyResponse, err
	}
	var termLocationMap search.TermLocationMap
	var ok bool
	if termLocationMap, ok = dm.Locations[fieldName]; !ok {
		return emptyResponse, errors.New(fmt.Sprintf("field '%s' missing", fieldName))
	}
	lines := make(map[uint]ResultEntry)
	for _, locations := range termLocationMap {
		// For every term location, either crate a new ResultEntry, or update existing ResultEntry with term locations
		for _, loc := range locations {
			// First, either look up a ResultEntry by its address or create a new ResultEntry
			var entry ResultEntry

			tmpLineAddr := uint(strings.LastIndexByte(body[:loc.Start], '\n') + 1)
			tmpLineEnd := strings.IndexRune(body[loc.End:], '\n')
			if tmpLineEnd != -1 {
				tmpLineEnd = int(loc.End) + tmpLineEnd
			} else {
				tmpLineEnd = len(body)
			}
			tmpSpanStart := uint(loc.Start) - tmpLineAddr
			tmpSpanEnd := uint(loc.End) - tmpLineAddr
			if mapVal, ok := lines[tmpLineAddr]; ok {
				// Adds the Span described by this `loc` to the existing ResultEntry.
				// Take care to avoid adding duplicate span in edge case
				entry = mapVal
				hasSpan := false
				for _, span := range entry.Spans {
					if span.Start == tmpSpanStart && span.End == tmpSpanEnd {
						hasSpan = true
						break
					}
				}
				if !hasSpan {
					entry.Spans = append(entry.Spans, Span{Start:tmpSpanStart, End:tmpSpanEnd})
				}
			} else {
				// Creates NewEntry with NoteID, Addr, Line, and the Span represented by this `loc`
				entry = ResultEntry{
					NoteID: dm.ID,
					Addr: tmpLineAddr,
					Line: body[tmpLineAddr:tmpLineEnd],
					Spans: []Span{{Start:tmpSpanStart, End:tmpSpanEnd}},
				}
			}
			lines[tmpLineAddr] = entry
		}
	}
	response := make([]ResultEntry, 0, len(lines))
	for _, entry := range lines {
		response = append(response, entry)
	}

	return response, nil
}

func mapSearchResult(fieldName string, searchResult *bleve.SearchResult) (AtlasResponse) {
	var ar AtlasResponse
	ar.ResultEntries = make([]ResultEntry, 0)
	for _, dm := range searchResult.Hits {
		resultEntries, err := mapDocumentMatchToResultEntrySlice(fieldName, dm)
		if err != nil {
			log.Println("cannot map DocumentMap, skipping", dm.ID, err)
			continue;
		}
		for _, entry := range resultEntries {
			ar.ResultEntries = append(ar.ResultEntries, entry)
		}
	}
	return ar
}

func resolveBody(fieldName string, dm *search.DocumentMatch) (string, error) {
	var value interface{}
	var ok bool
	if value, ok = dm.Fields[fieldName]; !ok {
		return "", errors.New(fmt.Sprintf("field '%s' missing", fieldName))
	}
	if s, ok := value.(string); ok {
		return s, nil
	}
	return fmt.Sprint(value), nil
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