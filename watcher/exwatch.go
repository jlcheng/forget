package watcher

import (
	"fmt"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/trace"
	rwatch "github.com/radovskyb/watcher"
	"io/ioutil"
	"os"
	"strings"
)

type WatcherFacade struct {
	Atlas *db.Atlas
}

// ReceiveWatchEvents will read events from thw watcher until the watcher is stopped. ReceiveWatchEvents will close
// the supplied stopCh when it is done.
func (wfacade *WatcherFacade) ReceiveWatchEvents(watcher *rwatch.Watcher, stopCh chan struct{}) {
	stop := false
	for !stop {
		select {
		case event := <-watcher.Event:
			trace.Debug(event)
			wfacade.onEvent(event)
		case err := <-watcher.Error:
			trace.Warn(err)
		case <-watcher.Closed:
			stop = true
		}
	}
	close(stopCh)
}

func (wfacade *WatcherFacade) onEvent(event rwatch.Event) {
	if event.IsDir() {
		return
	}

	if event.Op == rwatch.Chmod {
		return
	}

	if event.Op == rwatch.Remove {
		err := wfacade.Atlas.Remove(event.Path)
		if err != nil {
			trace.Warn("cannot remove: ", event.Path, err)
			return
		}
		trace.Debug("removed note: ", event.Path)
		return
	}

	trace.Debug("indexing file:", event.Path)
	info, err := os.Stat(event.Path)
	if err != nil {
		trace.Warn("cannot index:", event.Path, err)
		return
	}
	if strings.HasSuffix(event.Path, ".git") {
		trace.Debug("skipping %v", event.Path)
		return
	}

	// File name must contain a '.'
	if strings.LastIndexByte(event.Path, '.') < strings.LastIndexByte(event.Path, '/') {
		trace.Debug("ignoring file without dot: %v", rwatch.Event{}.Path)
		return
	}

	note := db.Note{
		ID: event.Path,
		Body: slurpFile(event.Path),
		Title: info.Name(),
		AccessTime: info.ModTime().Unix(),
	}
	wfacade.Atlas.Enqueue(note)
	trace.Debug("indexed file:", event.Path)
}

func slurpFile(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	defer f.Close()
	s, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return string(s)
}
