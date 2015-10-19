package main

import (
	"errors"
	"strconv"
	"strings"

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
	IsInt     bool
}

// ReplFlagValue holds a parsed ReplFlag
type ReplFlagValue struct {
	Flag      ReplFlag
	StringVal string
	BoolVal   bool
	IntVal    int
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

				parsedargs := []string{}

				for i := 0; i < len(args); i++ {
					// Check the flags
					arg := args[i]
					foundflag := false

					for _, flag := range cmd.Flags {
						if arg == "--"+flag.Name || arg == "-"+flag.ShortName {
							foundflag = true
							var newflagvalue ReplFlagValue

							// This is our flag!
							if flag.IsBool {
								newflagvalue = ReplFlagValue{Flag: flag, BoolVal: true}
							} else {
								// If parametered flag, get next arg
								i++
								if i < len(args) {
									argvalue := args[i]
									if flag.IsInt {
										newflagvalue = ReplFlagValue{Flag: flag, StringVal: argvalue}
									} else {
										val, _ := strconv.Atoi(argvalue)
										newflagvalue = ReplFlagValue{Flag: flag, IntVal: val}
									}
								} else {
									return nil, errors.New(arg + " need value")
								}
							}
							replYell.Flags = append(replYell.Flags, newflagvalue)
						}
					}
					if !foundflag {
						if strings.HasPrefix(arg, "-") {
							return nil, errors.New(arg + " unknown flag")
						}
						parsedargs = append(parsedargs, arg)
					}
				}
				// Now we parse for flags...
				//for _, arg := range args {
				// Do nothing for now... TODO!
				//}

				replYell.Args = parsedargs
				return &replYell, nil
				//cmd.Action(args, cliContext)
			}
		}
		return nil, errors.New(command + " not found")
	}

	return nil, nil
}

func replFlagIsSet(cmd *ReplYell, flagname string) bool {
	for _, flagvalue := range cmd.Flags {
		if flagvalue.Flag.Name == flagname {
			if flagvalue.Flag.IsBool {
				return flagvalue.BoolVal
			}
		}
	}
	return false
}

func replFlagIsInt(cmd *ReplYell, flagname string) int,
