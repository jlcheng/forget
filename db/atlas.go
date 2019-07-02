package db

import (
	"bytes"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/token/edgengram"
	"github.com/blevesearch/bleve/analysis/token/length"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search"
	"github.com/blevesearch/bleve/search/query"
	"github.com/jlcheng/forget/trace"
	"github.com/pkg/errors"
	"math"
	"strings"
)

const (
	BODY        = "Body"
	ACCESS_TIME = "AccessTime"
	TITLE       = "Title"

	DEFAULT_BATCH_SIZE = 1000
)

type Note struct {
	ID         string
	Body       string
	Title      string                  // some short title of this note
	Fragments  search.FieldFragmentMap // only used for query results, show a snippet of text around found terms
	AccessTime int64                   // time.Unix(), see NewIndexMapping():accessTime_fmap for FieldMapping
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

// Atlas is a lightweight interface over Bleve that batches indexing operations by default.
type Atlas struct {
	index bleve.Index  // Expected implementation is *bleve.indexImpl{}
	batch *bleve.Batch // supports batching
	size  int          // batch size
}

// Open attempts to reuse an existing index or create a new index at the path
func Open(path string, size int) (*Atlas, error) {
	runtimeConfig := map[string]interface{}{}
	index, err := bleve.OpenUsing(path, runtimeConfig)
	if err == nil {
		return NewAtlas(index, size), nil
	}
	switch err {
	case bleve.ErrorIndexMetaCorrupt, bleve.ErrorUnknownIndexType, bleve.ErrorUnknownStorageType:
		return nil, errors.Wrap(err, "cannot open existing index")
	case bleve.ErrorIndexPathDoesNotExist:
		// happy path
	default:
		return nil, errors.Wrap(err, "unexpected bleve error")
	}

	indexMapping, err := NewIndexMapping()
	if err != nil {
		return nil, err
	}
	index, err = bleve.NewUsing(path, indexMapping, scorch.Name, scorch.Name, nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open new index")

	}
	return NewAtlas(index, size), nil
}
func NewAtlas(index bleve.Index, size int) *Atlas {
	return &Atlas{
		index: index,
		batch: index.NewBatch(),
		size:  size,
	}
}

// Closes the Atlas instance and the underlying bleve index
func (s *Atlas) Close() error {
	err := s.Flush() // TODO: JCHENG handle returned error
	if err != nil {
		return err
	}
	err = s.index.Close()
	if err != nil {
		return errors.Wrap(err, "cannot close index")
	}
	return nil
}

func (s *Atlas) CloseQuietly() {
	err := s.Close()
	if err != nil {
		trace.Warn(err)
	}
}

func (s *Atlas) Enqueue(note Note) error {
	trace.Debug(fmt.Sprintf("Enqueue called for %v", note.Title))
	err := s.batch.Index(note.ID, note)
	if err != nil {
		return errors.WithStack(err)
	}
	if s.batch.Size() >= s.size {
		return s.Flush()
	}
	return nil
}

func (s *Atlas) Remove(noteID string) error {
	s.batch.Delete("")
	if s.batch.Size() >= s.size {
		return s.Flush()
	}
	return nil
}

func (s *Atlas) Flush() error {
	trace.Debug(fmt.Sprintf("Flush() called with batch.Size of %d", s.batch.Size()))
	err := s.index.Batch(s.batch)
	if err != nil {
		trace.Warn("Flush() failed")
		return errors.Wrap(err, "flush failed")
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
		return nil, errors.Wrap(err, "search failed")
	}
	notes := make([]Note, len(results.Hits))
	for idx := range notes {
		notes[idx] = toNote(results.Hits[idx])
	}
	return notes, nil
}

func (s *Atlas) rawIndex() bleve.Index {
	return s.index
}

func (s *Atlas) DumpAll() ([]Note, error) {
	if dc, err := s.index.DocCount(); err != nil {
		return nil, errors.Wrap(err, "error in DocCount")
	} else {
		trace.Debug("docCount", dc)
	}

	sr := bleve.NewSearchRequest(query.NewMatchAllQuery())
	sr.Fields = []string{"*"}
	results, err := s.index.Search(sr) // bleve/index_impl, bleve/search/collector/topn.Collect
	if err != nil {
		return nil, errors.Wrap(err, "search failed")
	}
	trace.Debug("hitsCount", len(results.Hits))

	notes := make([]Note, len(results.Hits))
	for idx := range notes {
		notes[idx] = toNote(results.Hits[idx])
	}
	return notes, nil
}

func NewIndexMapping() (mapping.IndexMapping, error) {
	indexMapping := bleve.NewIndexMapping()

	noteMapping := bleve.NewDocumentMapping()
	indexMapping.AddDocumentMapping("_default", noteMapping)

	const body_analyzer = "body_analyzer"
	bodyMapping := bleve.NewTextFieldMapping()
	noteMapping.AddFieldMappingsAt(BODY, bodyMapping)
	bodyMapping.Analyzer = body_analyzer

	accessTimeMapping := bleve.NewNumericFieldMapping()
	noteMapping.AddFieldMappingsAt(ACCESS_TIME, accessTimeMapping)

	const token_length_filter = "token_length_filter"
	err := indexMapping.AddCustomTokenFilter(token_length_filter,
		map[string]interface{}{
			"type": length.Name,
			"min":  2.0,
			"max":  32.0,
		})
	if err != nil {
		return nil, errors.Wrap(err, "cannot add custom token filter")
	}
	const edge_ngram_filter = "edge_ngram_filter"
	err = indexMapping.AddCustomTokenFilter(edge_ngram_filter,
		map[string]interface{}{
			"type": edgengram.Name,
			"min": 3.0,
			"max": 25.0,
		})
	if err != nil {
		return nil, errors.Wrap(err, "cannot add edgeNgram token filter")
	}

	err = indexMapping.AddCustomAnalyzer(body_analyzer, map[string]interface{}{
		"type":      custom.Name,
		"tokenizer": unicode.Name,
		"token_filters": []interface{}{
			lowercase.Name,
			en.StopName,
			token_length_filter,
			edge_ngram_filter,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot add custom analyzer")
	}

	return indexMapping, nil
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
