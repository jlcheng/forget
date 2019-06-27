package watcher

import (
	"github.com/jlcheng/forget/atlasrpc"
	"github.com/jlcheng/forget/testkit"
	"github.com/jlcheng/forget/trace"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

// TestRpcSearch is a Gateway Intergration Test
func TestRpcSearch(t *testing.T) {
	const RPC_PORT = 63999
	var err error
	baseTmpDir, err := ioutil.TempDir("", "4gt-index_")
	if err != nil {
		t.Error(err)
	}
	defer testkit.TempDirRemoveAll(baseTmpDir)
	indexDir := path.Join(baseTmpDir, "index")
	dataDir1 := path.Join(baseTmpDir, "data1")
	err = os.Mkdir(dataDir1, 0775)
	if err != nil {
		t.Error(err)
	}
	dataDir2 := path.Join(baseTmpDir, "data2")
	err = os.Mkdir(dataDir2, 0775)
	if err != nil {
		t.Error(err)
	}

	runexsvr := NewWatcherFacade()
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
	defer testkit.TempCloserClose(f)
	_, err = f.WriteString("this is my testcase")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Millisecond * 200)
	response, err := atlasrpc.RequestForBleveSearchResult("localhost", RPC_PORT, "testcase")
	if err != nil {
		t.Error(err)
	}
	if response.Total != 1 {
		t.Error("expected one result")
	}

}
