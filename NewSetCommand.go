package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"
)

// NewSetCommand configures the Set command
func NewSetCommand() ReplCommand {
	return ReplCommand{
		Name:  "set",
		Usage: "set key value",
		Flags: []ReplFlag{},
		Action: func(r *ReplYell) {
			setCommand(r)
		},
	}
}

func setCommand(r *ReplYell) {
	shellState := getShellState()

	if len(r.Args) != 2 {
		fmt.Println("Need a key and a value")
		return
	}

	key := r.Args[0]
	value := r.Args[1]

	if !strings.HasPrefix(key, "/") {
		if shellState.pwd != "/" {
			key = shellState.pwd + key
		} else {
			key = shellState.pwd + "/" + key
		}
		fmt.Println("Setting " + key)
	}

	resp, err := shellState.kapi.Set(context.TODO(), key, value, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	printResponseKey(resp)
}
