package atlasrpc

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	"time"
)

type SearchResult struct {
	Total uint64
	MaxScore float64
	Took time.Duration
	Hits search.DocumentMatchCollection
}

func BleveToAtlasSearchResult(src *bleve.SearchResult, dest *SearchResult) {
	dest.Total = src.Total
	dest.MaxScore = src.MaxScore
	dest.Took = src.Took
	dest.Hits = src.Hits
}
