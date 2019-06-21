package app

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	"github.com/jlcheng/forget/atlasrpc"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/trace"
	"github.com/jlcheng/forget/txtio"
	"github.com/pkg/errors"
	"strings"
	"time"
)

// GrepClient queries an Atlas server and renders the results similar to grep's output
func GrepClient(args []string) error {
	qterms := make([]string, len(args))
	for idx := range args {
		qterms[idx] = fmt.Sprintf(`+Body:"%s"`, args[idx])
	}
	stime := time.Now()

	sr, err := atlasrpc.RequestForBleveSearchResult(cli.Host(), cli.Port(), strings.Join(qterms, " "))
	if err != nil {
		return err
	}
	r := resultEntryList(sr)
	fmt.Printf("Found %v notes in %v\n", len(r), time.Since(stime))
	for _, entry := range r {
		fmt.Println(txtio.AnsiFmt(entry))
	}
	return nil
}

// resultEntryList converts a bleve.SearchResult to ResultEntries
func resultEntryList(sr *bleve.SearchResult) []db.ResultEntry {
	r := make([]db.ResultEntry, 0)
	for _, dm := range sr.Hits {
		chunk, err := mapDocumentMatchToResultEntrySlice("Body", dm)
		if err != nil {
			trace.Warn("cannot map DocumentMap, skipping", dm.ID, err)
			continue
		}
		r = append(r, chunk...)
	}
	return r
}

// mapDocumentMatchToResultEntrySlice flattens the DocumentMatch's
// Term=>Location map into a collection of matching lines. The error
// is non-nil if mapping fails.
func mapDocumentMatchToResultEntrySlice(fieldName string, dm *search.DocumentMatch) ([]db.ResultEntry, error) {
	emptyResponse := make([]db.ResultEntry, 0)
	body, err := resolveBody(fieldName, dm)
	if err != nil {
		return emptyResponse, err
	}
	var termLocationMap search.TermLocationMap
	var ok bool
	if termLocationMap, ok = dm.Locations[fieldName]; !ok {
		return emptyResponse, errors.Errorf("field '%s' missing", fieldName)
	}
	lines := make(map[uint]db.ResultEntry)
	for _, locations := range termLocationMap {
		// For every term location, either crate a new ResultEntry, or update existing ResultEntry with term locations
		for _, loc := range locations {
			// First, either look up a ResultEntry by its address or create a new ResultEntry
			var entry db.ResultEntry

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
					entry.Spans = append(entry.Spans, db.Span{Start: tmpSpanStart, End: tmpSpanEnd})
				}
			} else {
				// Creates NewEntry with NoteID, Addr, Line, and the Span represented by this `loc`
				entry = db.ResultEntry{
					NoteID: dm.ID,
					Addr:   tmpLineAddr,
					Line:   body[tmpLineAddr:tmpLineEnd],
					Spans:  []db.Span{{Start: tmpSpanStart, End: tmpSpanEnd}},
				}
			}
			lines[tmpLineAddr] = entry
		}
	}
	response := make([]db.ResultEntry, 0, len(lines))
	for _, entry := range lines {
		response = append(response, entry)
	}

	return response, nil
}

func resolveBody(fieldName string, dm *search.DocumentMatch) (string, error) {
	var value interface{}
	var ok bool
	if value, ok = dm.Fields[fieldName]; !ok {
		return "", errors.Errorf("field '%s' missing", fieldName)
	}
	if s, ok := value.(string); ok {
		return s, nil
	}
	return fmt.Sprint(value), nil
}
