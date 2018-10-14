package db

import "fmt"

type IndexBuffer struct {
	buf []Note
	notes chan(Note)
	stop chan(bool)

	maxBufSize int
	processed int
}

func (b IndexBuffer) StartProcessing() {
	go func() {
		fmt.Println("start")
		for {
			select {
			case n := <- b.notes:
				b.buf = append(b.buf, n)
				if len(b.buf) > b.maxBufSize {
					b.ProcessBuffer()
				}
			case v := <- b.stop:
				fmt.Println("stop received", v)
				//b.stop = nil
				break
			default:
				fmt.Println("default in select")
				break
			}
			fmt.Println("outside select")
		}
		fmt.Println("here")
		b.ProcessBuffer()
	}()
}

func (b IndexBuffer) ProcessBuffer() {
	//fmt.Println("buf.size", len(b.buf))
	for _, note := range b.buf {
		fmt.Println("indexing", note.ID)
		b.processed =+ 1
	}
	b.buf = b.buf[:0]
}

func (b IndexBuffer) SendNote(note Note) {
	select {
	case b.notes <- note:
		fmt.Println("note sent")
	default:
		fmt.Println("cannot send note")
	}
}