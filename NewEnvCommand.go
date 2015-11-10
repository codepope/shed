package main

import "fmt"

// NewEnvCommand configures the Set command
func NewEnvCommand() ReplCommand {
	return ReplCommand{
		Name:  "env",
		Usage: "change an shed setting",
		Flags: []ReplFlag{},
		Action: func(r *ReplYell) {
			envCommand(r)
		},
	}
}

func envCommand(r *ReplYell) {
	shellState := getShellState()
	if len(r.Args) == 0 {
		fmt.Println("Format:" + shellState.format)
		return
	}

	myvar := r.Args[0]
	switch myvar {
	case "s", "short":
		shellState.format = "short"
	case "j", "json":
		shellState.format = "json"
	case "p", "pretty":
		shellState.format = "pretty"
	}
	fmt.Println("Print:" + shellState.format)
}
