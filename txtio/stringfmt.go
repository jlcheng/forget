package txtio

import (
	"bytes"
	"github.com/jlcheng/forget/db"
	"sort"
)

func AnsiFmt(entry db.ResultEntry) string {
	const FCOLOR = "\033[96m"
	const RESET = "\033[0m"
	const TCOLOR = "\033[38;5;202;1m"

	// Writes out the colorized file name, followed by the line based on starting positions of Spans
	var buf bytes.Buffer
	progress := uint(0)  // Track amount of line already written

	buf.WriteString(FCOLOR)
	buf.WriteString(entry.NoteID)
	buf.WriteString(RESET)
	buf.WriteString(": ")

	// Ensure the Spans are sorted
	spans := make([]db.Span, len(entry.Spans))
	copy(spans, entry.Spans)
	sort.Slice(spans, func(i, j int) bool {
		return spans[i].Start < spans[j].Start
	})
	for _, span := range spans {
		if progress < span.Start {
			buf.WriteString(entry.Line[progress:span.Start])
		}
		buf.WriteString(TCOLOR)
		buf.WriteString(entry.Line[span.Start:span.End])
		buf.WriteString(RESET)
		progress = span.End
	}
	buf.WriteString(entry.Line[progress:])
	return buf.String()
}