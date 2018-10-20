package db

import (
	"fmt"
	"github.com/jlcheng/forget/testkit"
	"sync"
	"testing"
)

func TestNonFull(t *testing.T) {
	index, err := testkit.TempIndex()
	if err != nil {
		t.Error(err)
	}
	defer testkit.CleanUpIndex(t, index)

	b := NewBatcher(10, index)
	b.Send(Note{ID:"1"})
	b.Close()

	if count, _  := index.DocCount(); count != 1 {
		t.Error(count)
	}
}

func TestFull(t *testing.T) {
	index, err := testkit.TempIndex()
	if err != nil {
		t.Error(err)
	}
	defer testkit.CleanUpIndex(t, index)

	size := 10
	b := NewBatcher(uint(size), index)
	for i := 0; i < size; i++ {
		b.Send(Note{ID:fmt.Sprint(i)})
	}
	b.Close()

	if count, _  := index.DocCount(); int(count) != size {
		t.Error("count", count)
	}
}

func TestLargeSet(t *testing.T) {
	index, err := testkit.TempIndex()
	if err != nil {
		t.Error(err)
	}
	defer testkit.CleanUpIndex(t, index)

	size := 10000
	b := NewBatcher(uint(size), index)
	var wg sync.WaitGroup
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			b.Send(Note{ID: fmt.Sprint(idx)})
		}(i)
	}
	wg.Wait()
	b.Close()

	if count, _  := index.DocCount(); int(count) != size {
		t.Error("count", count)
	}
}