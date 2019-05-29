package rpc

import (
	"fmt"
	"github.com/jlcheng/forget/db"
	"github.com/pkg/errors"
	"log"
	"net"
	"net/rpc"
)

// ForgetService is an adapter around the Atlas instance which conforms to the net/rpc API.
type ForgetService struct {
	Atlas *db.Atlas
}

// Query exports atlas.QueryForResponse for net/rpc.
func (svc *ForgetService) Query(qstr string, reply *db.AtlasResponse) error {
	atlasResponse := svc.Atlas.QueryForResponse(qstr)
	reply.ResultEntries = atlasResponse.ResultEntries
	return nil
}

// QueryForBleveSearchResult exports atlas.QueryForBleveSearchResult for net/rpc.
func (svc *ForgetService) QueryForBleveSearchResult(qstr string, reply *Foo) error {
	searchResult, err := svc.Atlas.QueryForBleveSearchResult(qstr)
	if err != nil {
		return err
	}
	
	reply.Total = searchResult.Total
	return nil
}

// StartRpcServer starts a net/rpc server that exports the ForgetService. It blocks forever.
func StartRpcServer(atlas *db.Atlas, port int) {
	forgetService := ForgetService{Atlas: atlas}
	err := rpc.Register(&forgetService)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	rpc.Accept(l)
}

// Request makes a request to a ForgetService hosted at the specified host+port.
func Request(host string, port int, qstr string) (db.AtlasResponse, error) {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return db.AtlasResponse{}, err
	}
	var atlasResponse db.AtlasResponse
	err = client.Call("ForgetService.Query", qstr, &atlasResponse)
	if err != nil {
		return db.AtlasResponse{}, err
	}
	return atlasResponse, nil
}

type Foo struct {
	Total uint64
}

// RequestForBleveSearchResults makes a request to a ForgetService hosted at the specified host+port
func RequestForBleveSearchResult(host string, port int, qstr string) (*Foo, error) {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	foo := Foo{}
	err = client.Call("ForgetService.QueryForBleveSearchResult", qstr, &foo)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &foo, nil
}
