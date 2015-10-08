package main

import (
	"errors"

	"github.com/mattn/go-shellwords"
)

// ReplCommand - inspired by cli for repl commands
type ReplCommand struct {
	Name   string
	Usage  string
	Flags  []ReplFlag
	Action func(r *ReplYell)
}

// ReplFlag defines a flag where needed
type ReplFlag struct {
	Name      string
	ShortName string
	IsBool    bool
}

// ReplFlagValue holds a parsed ReplFlag
type ReplFlagValue struct {
	Flag      ReplFlag
	StringVal string
	BoolVal   bool
}

// ReplYell holds the completely parsed command
type ReplYell struct {
	Command ReplCommand
	Args    []string
	Flags   []ReplFlagValue
	Line    string
}

var parser *shellwords.Parser

func makeReplCommand(commands []ReplCommand, line string) (cmd *ReplYell, err error) {
	// Step through the command line
	// Find the command, find the ReplCommand
	// Process the string according to the ReplCommand

	// Lets use shellwords to break it up
	if parser == nil {
		parser = shellwords.NewParser()
		parser.ParseBacktick = false
		parser.ParseEnv = true
	}

	args, err := parser.Parse(line)
	if len(args) != 0 {
		command := args[0]
		args = args[1:]

		for _, cmd := range commands {
			if cmd.Name == command {
				// We have our command
				replYell := ReplYell{Command: cmd, Line: line}

				// Now we parse for flags...
				//for _, arg := range args {
				// Do nothing for now... TODO!
				//}

				replYell.Args = args
				return &replYell, nil
				//cmd.Action(args, cliContext)
			}
		}
		return nil, errors.New(command + " not found")
	}

	return nil, nil
}
