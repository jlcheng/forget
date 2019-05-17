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
	Warn(err)
	if trace, ok := err.(stackTracer); ok {
		for _, f := range trace.StackTrace() {
			fmt.Printf("%+s:%d\n", f, f)
		}
	}
}
