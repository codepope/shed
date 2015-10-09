package main

import (
	"fmt"
	"strings"

	"github.com/coreos/etcd/client"
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

		if !strings.HasPrefix(key, "/") {
			//Not absolute, add the present working directory
			argkey = shellState.pwd + argkey
		}

		// Is there a wildcard?
		if strings.Contains(argkey, "*") {
			// Find the last slash
			lastslash := strings.LastIndex(argkey, "/")

			if lastslash > 0 {
				// Split the string on it for the key
				key = argkey[0:lastslash]
				matcher = argkey
			} else {
				// Use the present working directory for the key
				key = shellState.pwd
				matcher = shellState.pwd + argkey
			}
		} else {
			key = argkey
		}
	} else {
		// No argument, so set the key to the present working dir
		key = shellState.pwd
	}

	if debug {
		fmt.Println("Key:" + key)
		fmt.Println("Matcher:" + matcher)
	}

	resp, err := shellState.kapi.Get(context.TODO(), key, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.Node.Dir {
		for _, node := range resp.Node.Nodes {
			if matcher == "" {
				printNode(node)
			} else if glob.Glob(matcher, node.Key) {
				printNode(node)
			}
		}
	} else {
		fmt.Println("Not a directory")
	}
}

func printNode(node *client.Node) {
	if node.Dir {
		fmt.Println(node.Key + "/")
	} else {
		fmt.Println(node.Key)
	}
}
