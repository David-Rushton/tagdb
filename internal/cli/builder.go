package cli

import (
	"os"
	"path"
)

// Builds a CLI command app.
type Builder struct {
	name        string
	version     string
	description string
	rootCommand *command
	exit        func(code int)
}

// Set the application name.
func (b *Builder) Name(name string) *Builder {
	b.name = name
	return b
}

// Set the application version.
func (b *Builder) Version(version string) *Builder {
	b.version = version
	return b
}

// The application description is used when returning help and usage tips.
func (b *Builder) Description(description string) *Builder {
	b.description = description
	return b
}

// Creates a new branch.
// Branches provide a way to group related commands together.
// If called directly a branch will display help for its direct subcommands.
//
// Example:
//
//	| binary | branch | command |
//	| ------ |--------| ------- |
//	| my_app | pod    | start   |
//	| my_app | pod    | stop    |
func (b *Builder) AddBranch(name, description string) (*command, error) {
	if b.rootCommand == nil {
		b.rootCommand = &command{
			name:        "",
			description: "command tree root",
			subcommands: map[string]*command{},
		}
	}

	invoke := newBranchCommandHandler(name, description)
	return addSubcommand(b.rootCommand, name, description, invoke)
}

// Adds a command.
func (b *Builder) AddCommand(name, description string, handler Invoker) (*command, error) {
	if b.rootCommand == nil {
		b.rootCommand = &command{
			name:        "",
			description: "command tree root",
			subcommands: map[string]*command{},
		}
	}

	return addSubcommand(b.rootCommand, name, description, handler)
}

func (b *Builder) Build() *app {
	if b.rootCommand == nil {
		b.rootCommand = &command{
			name:        "",
			description: "command tree root",
			subcommands: map[string]*command{},
		}
	}

	if b.exit == nil {
		b.exit = func(code int) {
			os.Exit(code)
		}
	}

	if b.name == "" {
		// Default to the executable name.
		_, file := path.Split(os.Args[1])
		b.name = file

	}

	return &app{
		name:        b.name,
		version:     b.version,
		description: b.description,
		rootCommand: b.rootCommand,
		exit:        b.exit,
	}
}
