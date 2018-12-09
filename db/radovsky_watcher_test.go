package db

import (
	"github.com/radovskyb/watcher"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"testing"
	"time"
)

// Used to learn the radovskyb/watcher API
func TestWatcherAPI(t *testing.T) {
	var tempDir string
	var err error
	if tempDir, err = ioutil.TempDir("", "forget-radovskyb-watcher"); err != nil {
		log.Fatal(err)
		return
	}

	events := make([]watcher.Event, 0, 0)

	w := watcher.New()
	var printGoRoutine sync.WaitGroup // notifies main thread that the following goroutine has stopped
	go func() {
		printGoRoutine.Add(1)
		stop := false
		for !stop {
			select {
			case event := <-w.Event:
				events = append(events, event)
			case err := <-w.Error:
				log.Fatal(err)
			case <-w.Closed:
				stop = true
			}
		}
		printGoRoutine.Done()
	}()

	if err := w.AddRecursive(tempDir); err != nil {
		log.Fatal(err)
	}

	go func() {
		w.Wait()
		file1 := path.Join(tempDir, "test1.txt")
		dir1 := path.Join(tempDir, "dir1")
		file2 := path.Join(tempDir, "dir1", "test2.txt")

		if err := ioutil.WriteFile(file1, []byte("test1"), 0644); err != nil {
			log.Fatal(err)
		}

		if err = os.Mkdir(dir1, 0755); err != nil {
			log.Fatal(err)
		} else {
			if err := ioutil.WriteFile(file2, []byte("test2"), 0644); err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(time.Millisecond * 1000)
		w.Close()
	}()

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatal(err)
	}

	printGoRoutine.Wait()

	wanted := watcher.Event{Op: watcher.Create, Path: path.Join(tempDir, "/test1.txt")}
	if !findEvent(events, wanted) {
		t.Fatalf("missing expected event: %v\n", wanted)
	}
	wanted = watcher.Event{Op: watcher.Create, Path: path.Join(tempDir, "/dir1/test2.txt")}
	if !findEvent(events, wanted) {
		t.Fatalf("missing expected event: %v\n", wanted)
	}
	wanted = watcher.Event{Op: watcher.Create, Path: path.Join(tempDir, "/dir1")}
	if !findEvent(events, wanted) {
		t.Fatalf("missing expected event: %v\n", wanted)
	}
	wanted = watcher.Event{Op: watcher.Create, Path: path.Join(tempDir, "/dir2")}
	if findEvent(events, wanted) {
		t.Fatalf("found unexpected event: %v\n", wanted)
	}

}

func findEvent(events []watcher.Event, target watcher.Event) bool {
	for _, event := range events {
		if event.Path == target.Path && event.Op == target.Op {
			return true
		}
	}
	return false
}
