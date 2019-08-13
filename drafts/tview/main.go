package main

import (
	"fmt"
	"github.com/jlcheng/forget/db"
	"github.com/rivo/tview"
)

// TviewUI has a search field and displays the search results
//
// Logically it contains two elements: A search input and a display
// for the results.
//
// The SearchInput exposes the current query and takes a callback
// named 'OnExecute'.
//
// The ResultDisplay renders a search result. The search result is
// pageable and contains snippet, relevance, and file path to each
// result. It receives a bleve.SearchResult object and renders it.
//
//
// type SearchResult struct {
// 	Status   *SearchStatus                  `json:"status"`
// 	Request  *SearchRequest                 `json:"request"`
// 	Hits     search.DocumentMatchCollection `json:"hits"`
// 	Total    uint64                         `json:"total_hits"`
// 	MaxScore float64                        `json:"max_score"`
// 	Took     time.Duration                  `json:"took"`
// 	Facets   search.FacetResults            `json:"facets"`
// }
//
type TviewUI struct {
	app *tview.Application
}

func NewTviewUI() *TviewUI {
	ui := new(TviewUI)
	ui.app = tview.NewApplication()
	return ui
}

func (ui *TviewUI) Run() {
	ui.app.Run()
}

func main() {
	foo := db.Atlas{}
	bar := NewTviewUI()
	bar.Run()
}
