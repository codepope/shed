package main

import (
	"fmt"
	"strings"

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

	if len(r.Args) <= 1 {
		var key string
		if len(r.Args) == 1 {
			key = r.Args[0]
		} else {
			key = ""
			fmt.Println("No key set: Getting pwd " + shellState.pwd)
		}

		if !strings.HasPrefix(key, "/") {
			key = shellState.pwd + key
		}

		resp, err := shellState.kapi.Get(context.TODO(), key, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		printResponseKey(resp)
	} else if len(r.Args) > 1 {
		fmt.Println("Need one key - multiple key retrieval TODO")
	}
}
