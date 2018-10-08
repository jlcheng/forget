package search

import (
	"github.com/jlcheng/forget/log"
	"time"
)

//go:generate echo hello world

type Document struct {
	Id string
	Body string
	AccessTime time.Time
}

type SearchEngine struct {

}

func (s SearchEngine) Enqueue(doc Document) error {
	log.Debug("enqueue for indexing:", doc.Id)
	return nil
}
