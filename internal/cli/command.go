package cli

import (
	"fmt"
)

// Represents a CLI command.
type command struct {
	name        string
	description string
	handler     commandHandler
	subcommands map[string]*command
}

// Represents a command handler function.
type commandHandler func() int

// Creates a top level command.
func (c *command) AddCommand(name, description string, handler commandHandler) (*command, error) {
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

func newBranchCommandHandler(name, description string) func() int {
	return func() int {
		fmt.Println(name)
		fmt.Println(description)
		fmt.Println()

		// TODO: Implement below.
		// for _, subcommand := range c.subcommands {
		// 	fmt.Println("`%s`:\t`%s`", c.name, c.description)

		// }
		// fmt.Println()

		return 0
	}
}

func addSubcommand(parent *command, name, description string, handler commandHandler) (*command, error) {
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
