package db

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	"github.com/blevesearch/bleve/search/query"
	"log"
	"sort"
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

type MatchContext struct {
	ContextLines []string  // The lines around the matched term
	MatchLineIndex uint    // The index in ContextLines where matched term(s) is found. This value is usually the
	                       // median of ContextLines, but for matches at the beginning and end of a document, this index
	                       // will be different, e.g., a matched term in the first line.
}

type MatchInfo struct {
	NoteID string             // The ID of the matched Note
	Score float64             // The relevance of the matched Note
	Contexts []*MatchContext  // Lines with matched terms and surrounding lines
}

// TODO: JCHENG estimate the typical size of this returned object. avg_line_length * 5 * avg_hits_per_doc * docs.
type QueryMatches struct {
	MatchInfoMap   map[string]*MatchInfo  // NoteID => *MatchInfo
}


// SRLocations - [S]earch [R]result Locations models where matched terms around found in a atlas.Note. It differs from
// the bleve data model by assuming a flat document model.
type SRLocations struct {
	NoteID string
	Body string
	Score float64
	TermLocationMap search.TermLocationMap
}

// QueryForMatches runs a query and returns a QueryMatches object
func (s *Atlas) QueryForMatches(qstr string) (QueryMatches, error) {
	queryMatches := QueryMatches{}

	q := query.NewQueryStringQuery(qstr)
	sr := bleve.NewSearchRequest(q)
	sr.SortBy([]string{ACCESS_TIME})
	sr.Fields = []string{"*"}
	sr.IncludeLocations = true
	searchResult, err := s.index.Search(sr)
	if err != nil {
		return queryMatches, err
	}

	queryMatches.MatchInfoMap = make(map[string]*MatchInfo)
	for _, dm := range searchResult.Hits {
		srLocations, ok := mapDocumentMatch("Body", dm)
		if ok {
			matchInfo := mapSRLocations(srLocations)
			queryMatches.MatchInfoMap[dm.ID] = &matchInfo
		}
	}

	return queryMatches, nil
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

// mapDocumentMatch maps a *DocumentMatch into a SRLocations. If the DocumentMatch object contains the named field
// then the returned SRLocation holds lines from the named field and the boolean is true. If the DocubmentMatch object
// does not contain the named field, then the SRLocation will be empty and the boolean will be false.
func mapDocumentMatch(fieldName string, dm *search.DocumentMatch) (SRLocations, bool) {
	var srLocations = SRLocations{
		NoteID: dm.ID,
		Score: dm.Score,
	}

	if body, ok := dm.Fields[fieldName]; !ok {
		return SRLocations{}, false
	} else {
		if s, ok := body.(string); !ok {
			srLocations.Body = fmt.Sprint(body)
		} else {
			srLocations.Body = s
 		}
	}

	if termLocationMap, ok := dm.Locations[fieldName]; !ok {
		return SRLocations{}, false
	} else {
		srLocations.TermLocationMap = termLocationMap
	}

	return srLocations, true
}

// Flattens the DocumentMatch's Term=>Location map into a collestction of matching lines. If, for some reason, we cannot
// transform this DocumentMatch into valid ResultEntry objects, a empty []ResultEntry will be returned and an non-nil
// error will be returned.
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

func mapSRLocations(srLocations SRLocations) MatchInfo {
	// For every matched term, copies its line into the MatchInfo.
	//
	// This method de-dupes lines, in case there are multiple matched terms on the same line.
	var matchInfo MatchInfo
	matchInfo.Contexts = make([]*MatchContext, 0)
	// Transform TermLocationMap into line-oriented structure
	lines := make(map[uint64]string)
	for _, locations := range srLocations.TermLocationMap {
		for _, location := range locations {
			// lineAddr is NOT the line number; Rather, it is the distance between the start of the text and the
			// first character of the line.
			lineAddr, line := getLineAround(srLocations.Body, location.Start, location.End)
			if _, ok := lines[lineAddr]; ok {
				continue
			}
			lines[lineAddr] = line
			matchContext := MatchContext{
				ContextLines: []string {line},
				MatchLineIndex: 0,
			}
			matchInfo.Contexts = append(matchInfo.Contexts, &matchContext)
		}
	}

	// Copy simple attributes over
	matchInfo.NoteID = srLocations.NoteID
	matchInfo.Score = srLocations.Score


	return matchInfo
}

/* == START: MatchQuery == */
func (qm *QueryMatches) NoteIDs() []string {
	var keys = make([]string, 0, len(qm.MatchInfoMap))
	for k := range qm.MatchInfoMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
/* == End: MatchQuery == */

/* == START: PrettyPrinter == */
func PPResultEntrySlice(resultEntries []ResultEntry) string {
	var buf bytes.Buffer
	buf.WriteString("[]ResultEntry:\n")
	for idx, entry := range resultEntries {
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