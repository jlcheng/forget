package db

import (
	"bytes"
	"fmt"
	"github.com/blevesearch/bleve/search"
)

func TermLocationToStr(termLocationMap *search.TermLocationMap) string {
	if termLocationMap == nil {
		return "nil"
	}

	var buf bytes.Buffer
	buf.WriteString("{\n")
	for k, v := range *termLocationMap {
		buf.WriteString(fmt.Sprintf("  %v: [", k))
		for idx, loc := range v {
			if idx != 0 {
				buf.WriteString(",")
			}
			buf.WriteString(LocationToStr(loc))
		}
		buf.WriteString("]")
	}

	buf.WriteString("\n}")

	return buf.String()
}

func LocationToStr(location *search.Location) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Location{Start:%v, End:%v, Pos:%v", location.Start, location.End, location.Pos))
	if len(location.ArrayPositions) != 0 {
		buf.WriteString(fmt.Sprintf(", ArrayPos:%v", location.ArrayPositions))
	}
	buf.WriteString("}")
	return buf.String()
}
