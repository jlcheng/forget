package search

import (
	"fmt"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/jlcheng/forget/debug"
	"time"
)

//go:generate echo hello world

type Document struct {
	Id string
	Body string
	AccessTime time.Time
}

type SearchEngine struct {
	index bleve.Index
}

func OpenIndex(path string) (*SearchEngine, error) {
	// func NewUsing(path string, mapping mapping.IndexMapping, indexType string, kvstore string, kvconfig map[string]interface{}) (Index, error) {
	index, err := bleve.NewUsing(path, NewIndexMapping(), scorch.Name, scorch.Name, nil)
	if err != nil {
		return nil, err
	}
	return &SearchEngine{
		index: index,
	}, nil
}

func (s *SearchEngine) Enqueue(doc Document) error {
	return s.index.Index(doc.Id, doc)
}

func (s *SearchEngine) Search(qstr string) ([]Document, error) {
	if dc, err := s.index.DocCount(); err != nil {
		return nil, err
	} else {
		debug.Debug("docCount", dc)
	}

	docs := make([]Document, 0)
	sr := bleve.NewSearchRequest(query.NewQueryStringQuery(qstr))
	sr.Fields = []string{"*"}
	results, err := s.index.Search(sr)
	if err != nil {
		return nil, err
	}
	for _, dm := range results.Hits {
		fmt.Printf("document: %v\n", dm.Document)
		fmt.Printf("document fields: %v\n", dm.Fields)
		fmt.Printf("document: %v\n", dm.String())
		doc := Document{
			Id: dm.ID,
			Body: dm.String(),
		}
		docs = append(docs, doc)
	}
	return docs, nil
}


func NewIndexMapping() mapping.IndexMapping {
	im := bleve.NewIndexMapping()
	im.DefaultAnalyzer = en.AnalyzerName
	dm := bleve.NewDocumentMapping()
	dm.DefaultAnalyzer = en.AnalyzerName
	bodyFieldMapping := bleve.NewTextFieldMapping()
	dm.AddFieldMappingsAt("Body", bodyFieldMapping)
	accessTimeFieldMapping := bleve.NewDateTimeFieldMapping()
	dm.AddFieldMappingsAt("AccessTime", accessTimeFieldMapping)
	im.DefaultMapping = dm

	return im
}