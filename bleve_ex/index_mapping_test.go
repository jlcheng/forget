package bleve

import (
	"errors"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/token/length"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/tokenizer/unicode"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/testkit"
	"testing"
)

func TestIndexMapping(t *testing.T) {
	const body_analyzer = "body_analyzer"
	const max_token_length = "max_token_length"

	noteMapping := bleve.NewDocumentMapping()
	noteMapping.DefaultAnalyzer = body_analyzer

	bodyMapping := bleve.NewTextFieldMapping()
	bodyMapping.Analyzer = body_analyzer

	noteMapping.AddFieldMappingsAt("Body", bodyMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("_default", noteMapping)

	var err error
	err = indexMapping.AddCustomTokenFilter(max_token_length,
		map[string]interface{}{
			"type": length.Name,
			"min":  5.0,
			"max":  7.0,
		})
	if err != nil {
		t.Fatal(err)
	}

	err = indexMapping.AddCustomAnalyzer(body_analyzer, map[string]interface{}{
		"type":      custom.Name,
		"tokenizer": unicode.Name,
		"token_filters": []interface{}{
			lowercase.Name,
			en.StopName,
			max_token_length,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = testkit.DoInTempIndexContext(func(index bleve.Index, tmpDir string) error {
		note := db.Note{
			ID:   "note1",
			Body: "a bc def four fives sixties sevenie eighties ninenines",
		}
		err := index.Index(note.ID, note)
		if err != nil {
			return err
		}
		if err = assertSearchResultLen(index, "a", 0); err != nil {
			return err
		}
		if err = assertSearchResultLen(index, "bc", 0); err != nil {
			return err
		}
		if err = assertSearchResultLen(index, "def", 0); err != nil {
			return err
		}
		if err = assertSearchResultLen(index, "four", 0); err != nil {
			return err
		}
		if err = assertSearchResultLen(index, "fives", 1); err != nil {
			return err
		}
		if err = assertSearchResultLen(index, "sixties", 1); err != nil {
			return err
		}
		if err = assertSearchResultLen(index, "eighties", 0); err != nil {
			return err
		}
		if err = assertSearchResultLen(index, "ninenines", 0); err != nil {
			return err
		}

		return nil
	}, indexMapping)
	if err != nil {
		t.Fatal(err)
	}
}

func assertSearchResultLen(index bleve.Index, query string, expected int) error {
	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return err
	}
	got := len(searchResult.Hits)
	if expected != got {
		return errors.New(fmt.Sprintf("expected: %d, got: %d, query: %s", expected, got, query))
	}
	return nil
}
