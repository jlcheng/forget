package db

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
)

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

// QueryForMatches runs a query and returns a QueryMatches object
func (s *Atlas) QueryForMatches(qstr string) QueryMatches {
	var searchResult *bleve.SearchResult
	_ = searchResult
	var hits search.DocumentMatchCollection = searchResult.Hits
	_ = hits
	var documentMatch *search.DocumentMatch = hits[0]
	_ = documentMatch
	var fieldTermLocationMap search.FieldTermLocationMap = documentMatch.Locations
	_ = fieldTermLocationMap
	// DocumentMatch contains
	//  - Locations, which is a FieldTermLocationMap (see later)
	//  - Fields, which is fieldName(string) => string

	// FieldTermLocationMap is a map of
	// fieldName(string) => TermLocationMap
	// Where TermLocationMap is a map of
	// term(string) => []*Locations
	//
	// Example:
	//  "Body" => {
	//    "Brown" => [ *Location(Start:0, End: 4), *Location(Start:10, End:14) ],
	//    "Dog"   => [ *Location(Start:5, End: 8), *Location(Start:15, End:18) ],
	//  }


	// As an example, assume the following text with a query of 'brown'
	//   the brown dog jumped over the red fox
	//   a brown bird flew over the red fox
	//   the brown chicken played
	//   with the red hen
	//   ...
	//   the end
	// yields...
	// TermLocationMap{
	//   brown: [Location{Start:4, End:9, Pos:2},Location{Start:40, End:45, Pos:10},Location{Start:77, End:82, Pos:18}]
	// }
	return QueryMatches{}
}

func documentMatchToMatchInfo(ff string, dm *search.DocumentMatch) *MatchInfo {
	// For every matching Location, place the line it is on into a MatchContext with only one line. Here, one must
	// de-dupe lines (in case there are multiple Location objects sharing the same line).
	// Then, each MatchContext is added to a map of NoteID(string)=>MatchInfo
	// The result of this operation procduces a mapping of NoteID=>MatchInfo, which allows one to enumerate matched
	// Notes and where (within each Note) the match occurs.
	//
	// Ignores any Location whose fieldName does not match `ff` (fieldName filter).
	matchInfo := MatchInfo{}
	matchInfo.Score = dm.Score
	for field, termLocationMap := range dm.Locations {
		if field != ff {
			continue
		}
		for term, locations := range termLocationMap {
			_ = term
			for _, location := range locations {
				mc := MatchContext{}
				mc.ContextLines = []string{fmt.Sprint(location)}
				matchInfo.Contexts = append(matchInfo.Contexts, &mc)
			}
		}
	}

	return &matchInfo
}