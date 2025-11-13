package cli

import "fmt"

type app struct {
	name        string
	version     string
	description string
	rootCommand *command
	exit        func(code int)
}

func (a *app) Run(args []string) {
	// Validation.
	if len(args) < 2 {
		fmt.Println("show help")
		a.exit(-1)
		return
	}

	// Get requested command.
	// Greedily consumes args until no matches are found.
	var candidateCommand *command
	var currentCommand = a.rootCommand

	argsQueue := newQueue(args[1:])

	for !argsQueue.isEmpty() {
		arg := argsQueue.peek().(string)

		if command, found := currentCommand.subcommands[arg]; found {
			currentCommand = command
			candidateCommand = command
			argsQueue.dequeue()
			continue
		}

		break
	}

	if candidateCommand == nil {
		// TODO: Print help for current command and exit 1.
		fmt.Println("cannot find command")
		a.exit(-1)
		return
	}

	// Execute command.
	fmt.Printf("executing `%s` with args `%s`\n", candidateCommand.name, argsQueue)
	fmt.Printf("`%s`\n\n", candidateCommand.description)
	a.exit(candidateCommand.handler())
}
