package subcmd

import (
	"github.com/jlcheng/forget/rpc"
	"github.com/jlcheng/forget/trace"
	"github.com/jlcheng/forget/watcher"
	"io/ioutil"
	"os"
	"path"
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
	go func() {
		err := runexsvr.Listen(RPC_PORT, indexDir, []string{dataDir1, dataDir2}, time.Millisecond*10)
		if err != nil {
			trace.PrintStackTrace(err)
			os.Exit(1)
		}
	}()
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
	if err != nil {
		t.Error(err)
	}
	if len(response.ResultEntries) != 1 {
		t.Fatal("expected one result")
	}
}
