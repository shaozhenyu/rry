package main

import (
	"os"

	"rry/client"
	"rry/server"

	"ember/cli"
)

func main() {
	cmds := cli.NewCmds()

	client.Reg(cmds)
	server.Reg(cmds)

	args := os.Args[1:]
	if len(args) == 0 {
		cmds.Help(true)
	} else {
		cmds.Run(args)
	}
}
