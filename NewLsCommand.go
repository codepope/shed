package main

import (
	"fmt"
	"strings"

	"github.com/ryanuber/go-glob"

	"golang.org/x/net/context"
)

// NewLsCommand configures the Set command
func NewLsCommand() ReplCommand {
	return ReplCommand{
		Name:  "ls",
		Usage: "list directory",
		Flags: []ReplFlag{},
		Action: func(r *ReplYell) {
			lsCommand(r)
		},
	}
}

func lsCommand(r *ReplYell) {
	shellState := getShellState()

	key := ""
	matcher := ""

	if len(r.Args) == 1 {
		argkey := r.Args[0]

		// Is there a wildcard?
		if strings.Contains(argkey, "*") {
			// Find the last slash
			lastslash := strings.LastIndex(argkey, "/")
			if lastslash != -1 {
				// Split the string on it for the key
				key = argkey[0:lastslash]
				matcher = argkey
			} else {
				// Use the present working directory for the key
				key = shellState.pwd
				matcher = shellState.pwd + argkey
			}
		} else {
			// No wild card. Is it absolute?
			if !strings.HasPrefix(key, "/") {
				//Not absolute, add the present working directory
				key = shellState.pwd + key
			}
		}
	} else {
		// No argument, so set the key to the present working dir
		key = shellState.pwd
	}

	resp, err := shellState.kapi.Get(context.TODO(), key, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.Node.Dir {
		for _, node := range resp.Node.Nodes {
			if matcher == "" {
				fmt.Println(node.Key)
			} else if glob.Glob(matcher, node.Key) {
				fmt.Println(node.Key)
			}
		}
	} else {
		fmt.Println("Not a directory")
	}

}
