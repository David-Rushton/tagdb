/*
Parses and dispatches command line arguments.

Supports commands, subcommands, arguments, and options.
*/
package cli

import (
	"fmt"
	"os"
	"sync"
)

type DuplicateCommandError struct {
	CommandName string
}

func (e *DuplicateCommandError) Error() string {
	return fmt.Sprintf("duplicate command `%s`", e.CommandName)
}

var (
	mutex    = sync.RWMutex{}
	commands = map[string]*command{}
)

type command struct {
	name        string
	description string
	handler     commandHandler
}

type commandHandler func() int

func NewCommand(name, description string, handler commandHandler) (*command, error) {
	// TODO: Validate input.

	mutex.Lock()
	defer mutex.Unlock()

	cmd := &command{
		name:        name,
		description: description,
		handler:     handler,
	}

	if _, found := commands[cmd.name]; found {
		return nil, &DuplicateCommandError{CommandName: cmd.name}
	}

	commands[cmd.name] = cmd

	return cmd, nil
}

func Run(args []string) {
	if len(args) < 2 {
		fmt.Println("show help")
		os.Exit(1)
	}

	if command, found := commands[args[1]]; found {
		os.Exit(command.handler())
	}

	fmt.Printf("cannot find command `%s`\n", args[0])
	os.Exit(1)
}
