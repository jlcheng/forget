package db

import (
	"github.com/blevesearch/bleve"
	"github.com/jlcheng/forget/trace"
	"sync"
)

type Batcher struct {
	buf        chan Note           // delivers incoming notes
	flush      chan bool           // signal that writer should flush writes

	doneSync   sync.WaitGroup      // allows for blocking Stop()
}

// NewBatcher returns a Batcher with the target batch size and callback.
// The bFunc callback will be invoked each time the buffer is filled or when the Batcher is stopped.
func NewBatcher(size uint, index bleve.Index) *Batcher {
	var b Batcher
	b.buf   = make(chan Note, size+1) // slightly larger than buffer size, so that the channel does not block
	b.flush = make(chan bool)
	b.doneSync.Add(1)
	go func() {
		// index  - bleve Index
		// batch  - current bleve.Batch
		// bcount - number of pending index ops
		// size   - if bcount >= size, commit
		defer b.doneSync.Done()
		batch  := index.NewBatch()
		bcount := uint(0)
		stop   := false
		for !stop {
			doFlush := false
			select {
			case note, open := <- b.buf:
				if !open {
					stop = true
					doFlush = true
					break
				}
				err := batch.Index(note.ID, note)
				if err != nil {
					trace.OnError(err)
				}
				bcount += 1
				if bcount >= size {
					doFlush = true
				}
			case <- b.flush:
				doFlush = true
			}
			if doFlush {
				index.Batch(batch)
				batch.Reset()
				bcount = 0
			}
		}
	}()
	return &b
}

func (b *Batcher) Send(note Note) {
	b.buf <- note
}

func (b *Batcher) Flush() {
	b.flush <- true
}

// Close waits for the receiver goroutine to finishonly call this when you're certain all senders have stopped sending.
func (b *Batcher) Close() {
	close(b.buf)
	b.doneSync.Wait()
}

