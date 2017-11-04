package server

import (
	"fmt"
	"sync"

	"ember/cli"
	"ember/http/rpc"
)

type Server struct {
	lock sync.Mutex
}

type Client struct {
	Gets func(path string) (err error)               `args:"path" return:"err"`
	Test func(path string) (value string, err error) `args:"path" return:"value,err"`
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

func (p *Server) Gets(path string) (err error) {
	fmt.Println("server get :", path)
	return
}

func (p *Server) Test(path string) (string, error) {
	fmt.Println("server test:", path)
	return path + "11", nil
}

func NewServer(path string) (*Server, error) {
	return &Server{}, nil
}
