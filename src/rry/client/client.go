package client

import (
	"fmt"
	"os"

	"ember/cli"
)

const (
	RemoteAddress = "127.0.0.1:8080"
	ConfigDir     = ".rry"
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

	err = InitLocalConfig(pwd, ConfigDir, username, args[0])
	cli.Check(err)

	err = InitRemoteConfig(pwd, ConfigDir, args[0])
	cli.Check(err)

	fmt.Printf("Init %s ok!\n", args[0])
}

func CmdSync(args []string) {
	Sync(args)
}

func Sync(args []string) {
	if len(args) > 1 {
		fmt.Println("usage: <bin> sync filepath")
		os.Exit(1)
	}

	pwd, err := os.Getwd()
	cli.Check(err)
	if len(args) == 1 {
		pwd = args[0]
	}

	node, err := NewNode(pwd, ConfigDir)
	cli.Check(err)

	err = node.Sync(pwd)
	cli.Check(err)

}

func CmdUpload(args []string) {
	fmt.Println("upload:", args)

	err := Upload("data/test", "data/test")
	cli.Check(err)
}

func CmdDownload(args []string) {
	fmt.Println("download:", args)

	err := Download("data/test", "data/test")
	cli.Check(err)
}
