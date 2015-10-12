package main

import (
	"fmt"
	"path"
	"strings"

	"golang.org/x/net/context"
)

// NewCdCommand configures the Set command
func NewCdCommand() ReplCommand {
	return ReplCommand{
		Name:  "cd",
		Usage: "set current directory",
		Flags: []ReplFlag{},
		Action: func(r *ReplYell) {
			cdCommand(r)
		},
	}
}

func cdCommand(r *ReplYell) {
	shellState := getShellState()

	switch len(r.Args) {
	case 0:
		// No args
		shellState.pwd = "/"
		// No need to check as there'll always be a root
		return
	case 1:
		newpath := r.Args[0]

		// do we have wildcards?
		if strings.Index(newpath, "*") > 0 {
			fmt.Println("No wildcards allowed currently")
			return
		}

		// Is this new path relative or absolute
		if !path.IsAbs(newpath) {
			newpath = path.Join(shellState.pwd, newpath)
		} else {
			newpath = path.Clean(newpath)
		}

		resp, err := shellState.kapi.Get(context.TODO(), newpath, nil)

		if err != nil {
			fmt.Println(err)
			return
		}

		if !resp.Node.Dir {
			fmt.Println("Not a directory")
			return
		}
		shellState.pwd = newpath
		return
	default:
		fmt.Println("Too many arguments. 'cd' takes one path")
	}
}
