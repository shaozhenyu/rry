package client

import (
	"fmt"
	"os"

	"rry/server"

	"ember/cli"
	"ember/http/rpc"
)

const (
	RemoteAddress = "127.0.0.1:8080"
	SzyDir        = ".szy"
	Sep           = "/"
)

func Reg(cmds *cli.Cmds) {
	cmds.Reg("init", "init current dir", CmdInit)
	cmds.Reg("sync", "sync file or dir", CmdSync)
	cmds.Reg("upload", "upload file or dir", CmdUpload)
	cmds.Reg("download", "download file or dir", CmdDownload)
}

func CmdInit(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: <bin> init path [username]")
		os.Exit(1)
	}

	username := "nobody"
	if len(args) > 1 {
		username = args[1]
	}

	pwd, err := os.Getwd()
	cli.Check(err)
	InitLocalConfig(pwd, SzyDir, username, args[0])
}

func CmdSync(args []string) {
	fmt.Println("sync:", args)
}

func CmdUpload(args []string) {
	fmt.Println("upload:", args)

	client := &server.Client{}
	rpc := rpc.NewClient(RemoteAddress)
	err := rpc.Reg(client)
	cli.Check(err)

	err = client.Gets("aa")
	cli.Check(err)

	err = Upload("data/test")
	cli.Check(err)
}

func CmdDownload(args []string) {
	fmt.Println("download:", args)

	err := Download("data/test")
	cli.Check(err)
}
