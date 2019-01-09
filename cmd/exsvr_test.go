package cmd

import (
	"github.com/jlcheng/forget/rpc"
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

	closeCh := make(chan struct{})
	defer close(closeCh)
	go RunExsvr(RPC_PORT, indexDir, []string{dataDir1, dataDir2}, closeCh)
	time.Sleep(time.Millisecond * 500)
	f, err := os.Create(path.Join(dataDir2, "first.txt"))
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	f.WriteString("this is my testcase")
	f.Close()
	time.Sleep(time.Second * 1)

	response, err := rpc.Request("localhost", RPC_PORT, "testcase")
	if expected := "found 1 note"; !strings.Contains(response, expected) {
		t.Error("missing substring: ", expected)
	}
}