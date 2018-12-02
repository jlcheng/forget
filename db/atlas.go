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

	DEFAULT_BATCH_SIZE = 1000
)

// The bleve-type resolves to "_default", see bleve/mapping/index.IndexMappingImpl.determineType()
type Note struct {
	ID         string
	Body       string
	Title      string                   // some short title of this note
	Fragments  search.FieldFragmentMap  // only used for query results, show a snippet of text around found terms
	AccessTime int64                    // time.Unix(), see NewIndexMapping():accessTime_fmap for FieldMapping
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
	// Expected implementation is *bleve.indexImpl{}
	index bleve.Index
	batch *bleve.Batch // supports batching
	size  int          // batch size
}

func Open(path string, size int) (*Atlas, error) {
	index, err := bleve.OpenUsing(path, map[string]interface{}{})
	if err == nil {
		return NewAtlas(index, size), nil
	}

	index, err = bleve.NewUsing(path, NewIndexMapping(), scorch.Name, scorch.Name, nil)
	if err != nil {
		return nil, err
	}
	return NewAtlas(index, size), nil
}
func NewAtlas(index bleve.Index, size int) *Atlas {
	return &Atlas{
		index: index,
		batch: index.NewBatch(),
		size: size,
	}
}

func (s *Atlas) Close() error {
	_ = s.Flush() // TODO: JCHENG handle returned error
	return s.index.Close()
}

func (s *Atlas) Enqueue(note Note) error {
	err := s.batch.Index(note.ID, note)
	if err != nil {
		return err
	}
	if s.batch.Size() > s.size {
		return s.Flush()
	}
	return nil
}

func (s *Atlas) Flush() error {
	err := s.index.Batch(s.batch)
	if err != nil {
		return err
	}
	s.batch.Reset()
	return nil
}

func (s *Atlas) GetDocCount() (uint64, error) {
	return s.index.DocCount()
}

func (s *Atlas) QueryString(qstr string) ([]Note, error) {
	q := query.NewQueryStringQuery(qstr)
	sr := bleve.NewSearchRequest(q)
	sr.SortBy([]string{ACCESS_TIME})
	sr.Fields = []string{"*"}
	sr.IncludeLocations = true
	results, err := s.index.Search(sr)
	if err != nil {
		return nil, err
	}
	notes := make([]Note, len(results.Hits))
	for idx, _ := range notes {
		notes[idx] = toNote(results.Hits[idx])
	}
	return notes, nil
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
		notes[idx] = toNote(results.Hits[idx])
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
	if doc.Fragments != nil {
		note.Fragments = doc.Fragments
	}
	return note
}