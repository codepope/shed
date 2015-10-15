package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

// NewLsCommand configures the Set command
func NewLsCommand() ReplCommand {
	return ReplCommand{
		Name:  "ls",
		Usage: "list directory",
		Flags: []ReplFlag{
			ReplFlag{Name: "recur", ShortName: "r", IsBool: true},
		},
		Action: func(r *ReplYell) {
			lsCommand(r)
		},
	}
}

func lsCommand(r *ReplYell) {
	shellState := getShellState()

	dir := ""
	matcher := ""

	switch len(r.Args) {

	case 0:
		dir = shellState.pwd
		matcher = path.Join(dir, "*")
	case 1:
		if len(r.Args) == 1 {
			argkey := r.Args[0]

			if !path.IsAbs(argkey) {
				argkey = path.Join(shellState.pwd, argkey)
			} else {
				argkey = path.Clean(argkey)
			}

			dir, _ = path.Split(argkey)
			matcher = argkey
		}
	}

	if *debug {
		fmt.Println("Key:" + dir)
		fmt.Println("Matcher:" + matcher)
	}

	options := &client.GetOptions{Recursive: replFlagIsSet(r, "recur")}

	resp, err := shellState.kapi.Get(context.TODO(), dir, options)

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.Node.Dir {
		for _, node := range resp.Node.Nodes {
			if matcher == "" {
				printNode(shellState, node)
			} else {
				match, err := path.Match(matcher, node.Key)
				if err != nil {
					fmt.Println(err)
					return // A swift exit
				}
				if match {
					printNode(shellState, node)
				}

			}
		}
	} else {
		fmt.Println("Not a directory")
	}
}

func printNode(shellState *ShellState, node *client.Node) {
	nodename := strings.TrimPrefix(node.Key, shellState.pwd+"/")
	if node.Dir {
		fmt.Println(nodename + "/")
		for _, subnode := range node.Nodes {
			printNode(shellState, subnode)
		}
	} else {
		fmt.Println(nodename)
	}
}
