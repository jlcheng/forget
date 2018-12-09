package watcher

import (
	"fmt"
	rwatch "github.com/radovskyb/watcher"
)


// ReceiveWatchEvents will read events from thw watcher until the watcher is stopped. ReceiveWatchEvents will close
// the supplied stopCh when it is done.
func ReceiveWatchEvents(watcher *rwatch.Watcher, stopCh chan struct{}) {
	stop := false
	for !stop {
		select {
		case event := <-watcher.Event:
			fmt.Println(event)
		case err := <-watcher.Error:
			fmt.Println(err)
		case <-watcher.Closed:
			stop = true
		}
	}
	close(stopCh)
}