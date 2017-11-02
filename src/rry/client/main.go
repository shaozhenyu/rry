package client

import (
	"fmt"

	"ember/cli"
)

const (
	SlSyncDir = ".slsync"
)

func Reg(cmds *cli.Cmds) {
	cmds.Reg("init", "init current dir", CmdInit)
	cmds.Reg("sync", "sync file or dir", CmdSync)
}

func CmdInit(args []string) {
	fmt.Println("init:", args)
}

func CmdSync(args []string) {
	fmt.Println("sync:", args)
}
