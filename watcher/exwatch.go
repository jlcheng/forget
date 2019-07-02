package watcher

import (
	"fmt"
	"github.com/jlcheng/forget/atlasrpc"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/db/files"
	"github.com/jlcheng/forget/trace"
	rwatch "github.com/radovskyb/watcher"
	"log"
	"strings"
	"time"
)

type WatcherFacade struct {
	watcher *rwatch.Watcher
}

// NewWatcherFacade creates a NewWatcherFacade instance
//
// The WatcherFacade object must be explicitly started with Listen()
func NewWatcherFacade() WatcherFacade {
	return WatcherFacade{
		watcher: rwatch.New(),
	}
}

// Listen starts a search server at the specified port.
//
// The search server accepts connections using Go RPC. It polls the dataDirs directories for changes based
// on the duration parameter. Changes and new files will be added to the index at indexDir.
func (wfacade *WatcherFacade) Listen(port int, indexDir string, dataDirs []string, duration time.Duration) error {
	trace.Debug(fmt.Sprintf("indexDir: %s", indexDir))
	trace.Debug(fmt.Sprintf("dataDirs: %s", strings.Join(dataDirs, ", ")))
	atlas, err := db.Open(indexDir, 1)
	if err != nil {
		return err
	}
	defer atlas.CloseQuietly()
	docCount, err := atlas.GetDocCount()
	if err != nil {
		return err
	}
	trace.Debug("atlas doc count:", docCount)

	fmt.Printf("Starting rpc on port %d\n", port)
	go atlasrpc.StartRpcServer(atlas, port)

	// Creates a radovskyb.Watcher. Starts listening to its events. Finally, start the Watcher.
	for _, dataDir := range dataDirs {
		if err := wfacade.watcher.AddRecursive(dataDir); err != nil {
			log.Fatalf("cannot watch %v: %v\n", dataDir, err)
		}
	}

	go ReceiveWatchEvents(atlas, wfacade.watcher)

	return wfacade.watcher.Start(duration)
}

func (wfacade *WatcherFacade) Close() {
	wfacade.watcher.Close()
}

// ReceiveWatchEvents will delegate relevant filesystem events to an db.Atlas instance.
func ReceiveWatchEvents(atlas *db.Atlas, watcher *rwatch.Watcher) {
	stop := false
	for !stop {
		select {
		case event := <-watcher.Event:
			onEvent(atlas, event)
		case err := <-watcher.Error:
			trace.Warn(err)
		case <-watcher.Closed:
			stop = true
		}
	}
}

func onEvent(atlas *db.Atlas, event rwatch.Event) {
	if event.IsDir() {
		return
	}
	path := event.Path

	switch event.Op {
	case rwatch.Chmod:
		trace.Debug(event)
	case rwatch.Remove:
		err := atlas.Remove(path)
		if err != nil {
			trace.Warn("cannot remove: ", path, err)
		} else {
			trace.Debug("removed note: ", path)
		}
	case rwatch.Create, rwatch.Write:
		if !db.FilterFile(path, event.FileInfo) {
			trace.Debug(fmt.Sprintf("no-index: [%v] %v", path, event.Op))
		} else {
			notes, err := files.ParseFile(path)
			if err != nil {
				trace.Warn(fmt.Sprintf("cannot index [%v]: %v", path, err))
			}
			for _, note := range notes {
				err = atlas.Enqueue(note)
				if err != nil {
					trace.Warn(fmt.Sprintf("cannot index [%v]: %v", note.ID, err))
				}
			}

		}
	case rwatch.Rename, rwatch.Move:
		trace.Warn(fmt.Sprintf("not-implemented: %v, %v", event.Op, path))
	default:
		trace.Warn(fmt.Sprintf("not-implemented: %v", event))
	}
}
