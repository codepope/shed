package main

import "fmt"

// NewPwdCommand configures the Set command
func NewPwdCommand() ReplCommand {
	return ReplCommand{
		Name:  "pwd",
		Usage: "show current directory",
		Flags: []ReplFlag{},
		Action: func(r *ReplYell) {
			pwdCommand(r)
		},
	}
}

func pwdCommand(r *ReplYell) {
	shellState := getShellState()
	fmt.Println(shellState.pwd)
}
