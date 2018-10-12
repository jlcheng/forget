package recall

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search/query"
	"github.com/jlcheng/forget/trace"
)

//go:generate echo hello world

// The bleve-type resolves to "_default", see bleve/mapping/index.IndexMappingImpl.determineType()
type Document struct {
	ID         string
	Body       string
	AccessTime int64    // time.Unix(), see NewIndexMapping():accessTime_fmap for FieldMapping
}

type Atlas struct {
	// here's to hope that bleve+scorch goes the way of lucene rather than kestrel
	// expected impl is blevesearch/bleve.indexImpl
	index bleve.Index
}

func Open(path string) (*Atlas, error) {
	index, err := bleve.NewUsing(path, NewIndexMapping(), scorch.Name, scorch.Name, nil)
	if err != nil {
		return nil, err
	}
	return &Atlas{
		index: index,
	}, nil
}

func (s *Atlas) Enqueue(doc Document) error {
	return s.index.Index(doc.ID, doc)
}

func (s *Atlas) Search(qstr string) ([]Document, error) {
	if dc, err := s.index.DocCount(); err != nil {
		return nil, err
	} else {
		trace.Debug("docCount", dc)
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
			ID:   dm.ID,
			Body: dm.String(),
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func NewIndexMapping() mapping.IndexMapping {
	imap := bleve.NewIndexMapping()

	// needed because bleve will map atlas.Document to the "_default" bleve-type
	main_dmap := bleve.NewDocumentMapping()
	imap.AddDocumentMapping("_default", main_dmap)

	// configure the fields in atlas.Document, excepting doc.ID
	body_fmap := bleve.NewTextFieldMapping()
	main_dmap.AddFieldMappingsAt("Body", body_fmap)
	accessTime_fmap := bleve.NewNumericFieldMapping()
	main_dmap.AddFieldMappingsAt("AccessTime", accessTime_fmap)

	return imap
}

func NewDefaultDocumentMapping() mapping.DocumentMapping {

}