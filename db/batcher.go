package db

import "github.com/blevesearch/bleve"

type BBatcher struct {
	index bleve.Index  // performs indexing
	batch *bleve.Batch // supports batching
	size  int          // batch size
}

func NewBBatcher(index bleve.Index, size int) *BBatcher {
	return &BBatcher{
		index: index,
		batch: index.NewBatch(),
		size:  size,
	}
}

func (b *BBatcher) Index(note Note) error {
	err := b.batch.Index(note.ID, note)
	if err != nil {
		return err
	}
	if b.batch.Size() > b.size {
		return b.Flush()
	}
	return nil
}

func (b *BBatcher) Flush() error {
	err := b.index.Batch(b.batch)
	if err != nil {
		return err
	}
	b.batch.Reset()
	return nil
}
