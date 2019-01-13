package watcher

import (
	"fmt"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/rpc"
	"github.com/jlcheng/forget/trace"
	rwatch "github.com/radovskyb/watcher"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type WatcherFacade struct {
	watcher *rwatch.Watcher
}

func NewWatcherFacade() WatcherFacade {
	return WatcherFacade{
		watcher: rwatch.New(),
	}
}

func (wfacade *WatcherFacade) Listen(port int, indexDir string, dataDirs []string, duration time.Duration) error {
	trace.Debug(fmt.Sprintf("indexDir: %s", indexDir))
	trace.Debug(fmt.Sprintf("dataDirs: %s", strings.Join(dataDirs, ", ")))
	atlas, err := db.Open(indexDir, 1)
	if err != nil {
		return err
	}
	docCount, err := atlas.GetDocCount()
	trace.Debug("atlas doc count:", docCount)

	fmt.Printf("Starting rpc on port %d\n", port)
	go rpc.StartRpcServer(atlas, port)

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
			trace.Debug(event)
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

	if event.Op == rwatch.Chmod {
		return
	}

	path := event.Path

	// Do not index any files under .git
	for tmpPath := path;
		tmpPath != "." && tmpPath != string(filepath.Separator);  // invariant
	tmpPath = filepath.Dir(tmpPath) {
		tmpInfo, err := os.Stat(tmpPath)
		if os.IsNotExist(err) {
			// Remove event, in which case the file does not exist but parent may be .git
			continue
		}
		if err != nil {
			trace.Warn("cannot index:", path, err)
			return
		}
		if tmpInfo.IsDir() && strings.HasSuffix(tmpPath, ".git") {
			trace.Debug(fmt.Sprintf("skipping %v under .git", path))
			return
		}
	}

	if event.Op == rwatch.Remove {
		err := atlas.Remove(path)
		if err != nil {
			trace.Warn("cannot remove: ", path, err)
			return
		}
		trace.Debug("removed note: ", path)
		return
	}

	// File name must contain a '.'
	if strings.LastIndexByte(path, '.') < strings.LastIndexByte(path, '/') {
		trace.Debug("ignoring file without dot: %v", rwatch.Event{}.Path)
		return
	}

	// Omitting large files
	const ONE_MB = int64(1024 * 1024)
	if event.FileInfo.Size() > ONE_MB {
		trace.Debug("skipping %v (too large)", path)
		return
	}

	note := db.Note{
		ID: path,
		Body: slurpFile(path),
		Title: event.FileInfo.Name(),
		AccessTime: event.FileInfo.ModTime().Unix(),
	}
	atlas.Enqueue(note)
	trace.Debug("indexed file:", path)
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
