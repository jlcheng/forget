package rpc

import (
	"bytes"
	"fmt"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/txtio"
	"log"
	"net"
	"net/rpc"
	"time"
)

type ForgetService struct {
	Atlas *db.Atlas
}

func (svc *ForgetService) Query(qstr string, reply *string) error {
	stime := time.Now()
	atlasResponse := svc.Atlas.QueryForResponse(qstr)
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("found %v notes in %v\n", len(atlasResponse.ResultEntries), time.Since(stime)))
	for _, entry := range atlasResponse.ResultEntries {
		buf.WriteString(fmt.Sprintln(txtio.AnsiFmt(entry)))
	}

	*reply = buf.String()
	return nil
}

func StartRpcServer(atlas *db.Atlas, port int) {
	forgetService := ForgetService{Atlas: atlas}
	rpc.Register(&forgetService)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	rpc.Accept(l)
}

func Request(host string, port int, qstr string) (string, error) {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return "", err
	}
	var response string
	err = client.Call("ForgetService.Query", qstr, &response)
	if err != nil {
		return "", err
	}
	return response, nil
}