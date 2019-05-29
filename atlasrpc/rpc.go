package atlasrpc

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/jlcheng/forget/db"
	"github.com/pkg/errors"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// ForgetService is an adapter around the Atlas instance which conforms to the net/rpc API.
type ForgetService struct {
	Atlas *db.Atlas
}

// QueryForBleveSearchResult exports atlas.QueryForBleveSearchResult for net/rpc.
func (svc ForgetService) QueryForBleveSearchResult(qstr string, reply *bleve.SearchResult) error {
	searchResult, err := svc.Atlas.QueryForBleveSearchResult(qstr)
	if err != nil {
		return err
	}

	*reply = *searchResult
	return nil
}

// StartRpcServer starts a net/rpc server that exports the ForgetService. It blocks forever.
func StartRpcServer(atlas *db.Atlas, port int) {
	forgetService := ForgetService{Atlas: atlas}
	server := rpc.NewServer()
	err := server.Register(forgetService)
	if err != nil {
		log.Fatal("here")
	}
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		if conn, err := l.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}

// Request makes a request to a ForgetService hosted at the specified host+port.
func Request(host string, port int, qstr string) (db.AtlasResponse, error) {
	client, err := jsonrpc.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
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

// RequestForBleveSearchResults makes a request to a ForgetService hosted at the specified host+port
func RequestForBleveSearchResult(host string, port int, qstr string) (*bleve.SearchResult, error) {
	client, err := jsonrpc.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var sr bleve.SearchResult
	err = client.Call("ForgetService.QueryForBleveSearchResult", qstr, &sr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &sr, nil
}
