package db

import (
	"testing"
	"time"
)

var tBuf IndexBuffer

func reset() {
	tBuf = IndexBuffer{
		buf: make([]Note, 0),
		notes: make(chan Note, 64),
		stop: make(chan bool),
		maxBufSize: 64,
		processed: 0,
	}
}

func TestSend(t *testing.T) {
	reset()
	tBuf.StartProcessing()
	tBuf.SendNote(NewNote("1"))
	close(tBuf.stop)
	time.Sleep(4*time.Second)
}

func NewNote(ID string) Note {
	return Note{
		ID: ID,
		Body: "",
		Title: "",
		Fragments: nil,
		AccessTime: 0,
	}
}
