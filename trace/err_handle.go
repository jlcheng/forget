package trace

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
)

func TryClose(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		Warn(err)
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func PrintStackTrace(err error) {
	if _, ok := err.(stackTracer); ok {
		fmt.Printf("%+v\n", err)
	} else {
		fmt.Println(err)
	}
}
