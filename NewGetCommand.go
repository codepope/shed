package main

import (
	"fmt"
	"path"

	"golang.org/x/net/context"
)

// NewGetCommand configures the Set command
func NewGetCommand() ReplCommand {
	return ReplCommand{
		Name:  "get",
		Usage: "get key value",
		Flags: []ReplFlag{},
		Action: func(r *ReplYell) {
			getCommand(r)
		},
	}
}

func getCommand(r *ReplYell) {
	shellState := getShellState()
	key := ""

	switch len(r.Args) {
	case 1:
		key = r.Args[0]
	case 0:
		key = ""
		fmt.Println("No key set: Getting pwd " + shellState.pwd)
	default:
		fmt.Println("Need one key - multiple key retrieval TODO")
		return
	}

	if !path.IsAbs(key) {
		key = path.Join(shellState.pwd, key)
	} else {
		key = path.Clean(key)
	}

	resp, err := shellState.kapi.Get(context.TODO(), key, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	printResponseKey(resp)

}
