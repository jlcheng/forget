package db

import (
	"bytes"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search"
	"github.com/blevesearch/bleve/search/query"
	"github.com/jlcheng/forget/trace"
	"math"
	"strings"
)

//go:generate echo hello world

const (
	BODY = "Body"
	ACCESS_TIME = "AccessTime"
	TITLE = "Title"
)

// The bleve-type resolves to "_default", see bleve/mapping/index.IndexMappingImpl.determineType()
type Note struct {
	ID         string
	Body       string
	Title      string   // some short title of this note
	AccessTime int64    // time.Unix(), see NewIndexMapping():accessTime_fmap for FieldMapping
}

func (s Note) PrettyString() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("note.ID: %v\n", s.ID))
	buf.WriteString(fmt.Sprintf("  note.AccessTime: %d\n", s.AccessTime))
	buf.WriteString(fmt.Sprintf("  note.Title: %v\n", s.Title))
	snippet := s.Body
	snippet = strings.Replace(snippet, "\n", " ", -1)
	snippet = snippet[:int(math.Min(80, float64(len(snippet))))]
	buf.WriteString(fmt.Sprintf("  note.Body: %v\n", snippet))
	return buf.String()
}

type Atlas struct {
	// here's to hope that bleve+scorch goes the way of lucene rather than kestrel
	// expected impl is blevesearch/bleve.indexImpl
	index bleve.Index
}

func Open(path string) (*Atlas, error) {
	index, err := bleve.OpenUsing(path, map[string]interface{}{})
	if err == nil {
		return &Atlas{index:index}, nil
	}

	index, err = bleve.NewUsing(path, NewIndexMapping(), scorch.Name, scorch.Name, nil)
	if err != nil {
		return nil, err
	}
	return &Atlas{
		index: index,
	}, nil
}

func (s *Atlas) Close() error {
	return s.index.Close()
}

func (s *Atlas) Enqueue(doc Note) error {
	return s.index.Index(doc.ID, doc)
}

func (s *Atlas) DumpAll() ([]Note, error) {
	if dc, err := s.index.DocCount(); err != nil {
		return nil, err
	} else {
		trace.Debug("docCount", dc)
	}

	sr := bleve.NewSearchRequest(query.NewMatchAllQuery())
	sr.Fields = []string{"*"}
	results, err := s.index.Search(sr)  // bleve/index_impl, bleve/search/collector/topn.Collect
	if err != nil {
		return nil, err
	}
	trace.Debug("hitsCount", len(results.Hits))

	notes := make([]Note, len(results.Hits))
	for idx, _ := range notes {
		dm := results.Hits[idx]
		notes[idx] = toNote(dm)
	}
	return notes, nil
}

func NewIndexMapping() mapping.IndexMapping {
	imap := bleve.NewIndexMapping()

	// needed because bleve will map atlas.Note to the "_default" bleve-type
	main_dmap := bleve.NewDocumentMapping()
	imap.AddDocumentMapping("_default", main_dmap)

	// configure the fields in atlas.Note, excepting doc.ID - necessary?
	body_fmap := bleve.NewTextFieldMapping()
	main_dmap.AddFieldMappingsAt(BODY, body_fmap)
	accessTime_fmap := bleve.NewNumericFieldMapping()
	main_dmap.AddFieldMappingsAt(ACCESS_TIME, accessTime_fmap)

	return imap
}

func toNote(doc *search.DocumentMatch) Note {
	note := Note{}
	note.ID = doc.ID
	if atime, ok := doc.Fields[ACCESS_TIME]; ok {
		if v, ok := atime.(float64); ok {
			note.AccessTime = int64(v)
		}
	}
	if body, ok := doc.Fields[BODY]; ok {
		if v, ok := body.(string); ok {
			note.Body = v
		}
	}
	if title, ok := doc.Fields[TITLE]; ok {
		if v, ok := title.(string); ok {
			note.Title = v
		}
	}
	return note
}