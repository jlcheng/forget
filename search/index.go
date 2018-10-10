package search

import (
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve"
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
	return s.index.Index("id", doc)
}


func NewIndexMapping() mapping.IndexMapping {
	im := bleve.NewIndexMapping()
	dm := bleve.NewDocumentMapping()
	bodyFieldMapping := bleve.NewTextFieldMapping()
	dm.AddFieldMappingsAt("Body", bodyFieldMapping)
	accessTimeFieldMapping := bleve.NewDateTimeFieldMapping()
	dm.AddFieldMappingsAt("AccessTime", accessTimeFieldMapping)
	im.DefaultMapping = dm
	return im
}