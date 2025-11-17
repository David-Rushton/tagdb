package cli

import (
	"fmt"
)

type Invoker interface {
	Invoke() int
}

type invokeFunc func() int

func (f invokeFunc) Invoke() int {
	return f()
}

// Represents a CLI command.
type command struct {
	name        string
	description string
	handler     Invoker
	subcommands map[string]*command
}

// Creates a top level command.
func (c *command) AddCommand(name, description string, handler Invoker) (*command, error) {
	return addSubcommand(c, name, description, handler)
}

// Creates a new branch.
// Branches provide a way to group related commands together.
// If called directly a branch will display help for its direct subcommands.
//
// Example:
//
//	| binary | branch | command |
//	| ------ | ------ | ------- |
//	| my_app | pod    | start   |
//	| my_app | pod    | stop    |
func (c *command) AddBranch(name, description string) (*command, error) {
	return addSubcommand(c, name, description, newBranchCommandHandler(name, description))
}

func newBranchCommandHandler(name, description string) Invoker {
	invoke := func() int {
		fmt.Printf("Usage: %s [command]\n\n", name)
		fmt.Println(description)
		fmt.Println("\nCommands:")
		return 0
	}
	invokeFunc := invokeFunc(invoke)
	return invokeFunc
}

func addSubcommand(parent *command, name, description string, handler Invoker) (*command, error) {
	// Validation.
	if err := validateCommandName(name); err != nil {
		return nil, err
	}

	if err := validateDescription(description); err != nil {
		return nil, err
	}

	// Create command.
	cmd := &command{
		name:        name,
		description: description,
		handler:     handler,
		subcommands: map[string]*command{},
	}

	// Register command.
	if _, found := parent.subcommands[cmd.name]; found {
		return nil, &DuplicateCommandError{CommandName: cmd.name}
	}
	parent.subcommands[cmd.name] = cmd

	return cmd, nil
}
