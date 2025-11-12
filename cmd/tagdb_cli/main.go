package main

import (
	"fmt"
	"os"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/cli"
)

func main() {
	fmt.Println("cli parsing test")
	fmt.Println("----------------")
	fmt.Println()

	cli.NewCommand("good", "a command that should exit 0", goodHandler)
	cli.NewCommand("bad", "a command that should exit 1", badHandler)

	cli.Run(os.Args)
}

func goodHandler() int {
	fmt.Println("good command called")
	return 0
}

func badHandler() int {
	fmt.Println("bad command called")
	return 1
}
