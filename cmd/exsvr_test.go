package cmd

import (
	"github.com/jlcheng/forget/rpc"
	"github.com/jlcheng/forget/watcher"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func TestRunExsvr(t *testing.T) {
	const RPC_PORT = 63999
	var err error
	indexDir, err := ioutil.TempDir("", "4gt-index")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(indexDir)
	dataDir1, err := ioutil.TempDir("", "4gt-data1")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dataDir1)
	dataDir2, err := ioutil.TempDir("", "4gt-data2")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dataDir2)

	runexsvr := watcher.NewWatcherFacade()
	defer runexsvr.Close()
    go runexsvr.Listen(RPC_PORT, indexDir, []string{dataDir1, dataDir2}, time.Millisecond * 10)
	time.Sleep(time.Millisecond * 100)
	f, err := os.Create(path.Join(dataDir2, "first.txt"))
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	f.WriteString("this is my testcase")
	f.Close()
	time.Sleep(time.Millisecond * 200)
	response, err := rpc.Request("localhost", RPC_PORT, "testcase")
	if expected := "found 1 note"; !strings.Contains(response, expected) {
		t.Error("missing substring: ", expected)
	}
}