package server

import (
	"fmt"
	"sync"

	"rry/share"

	"ember/cli"
	"ember/http/rpc"
)

type Server struct {
	service *share.Service
	lock    sync.Mutex
}

type Client struct {
	Gets func(path string) (leaves share.Leaves, err error)         `args:"path" return:"leaves,err"`
	Put  func(path string, value string, clocks string) (err error) `args:"path,value,clocks"`
}

func Reg(cmds *cli.Cmds) {
	cmds.Reg("run", "run slsync server", CmdRun)
}

func CmdRun(args []string) {

	dir := "data"
	sobj, err := NewServer(dir)
	cli.Check(err)
	rpc := rpc.NewServer()
	err = rpc.Reg(sobj, &Client{})
	cli.Check(err)
	err = rpc.Run("/", 8080)
	cli.Check(err)
}

func (p *Server) Gets(path string) (leaves share.Leaves, err error) {
	fmt.Println("gets : ", path)
	leaves = p.service.GetAll(path)
	return
}

func (p *Server) Put(path string, clocks string, value string) (err error) {
	cs, err := share.NewClocksFromString(clocks)
	if err != nil {
		return
	}
	err = p.service.Put(path, cs, value)
	return
}

func NewServer(path string) (*Server, error) {
	service, err := share.NewService(path)
	if err != nil {
		return nil, err
	}
	return &Server{service: service}, nil
}
