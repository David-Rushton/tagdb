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

	builder := cli.Builder{}
	builder.Name("tagdb-cli")
	builder.Version("0.1.0-test")
	builder.Description("A CLI for TagDB")

	branch, err := builder.AddBranch("wip", "testing api structure")
	if err != nil {
		panic(err)
	}

	_, err = branch.AddCommand("good", "a good command", goodHandler)
	if err != nil {
		panic(err)
	}

	_, err = branch.AddCommand("bad", "a bad command", badHandler)
	if err != nil {
		panic(err)
	}

	app := builder.Build()
	app.Run(os.Args)
}

func goodHandler() int {
	fmt.Println("good command called")
	return 0
}

func badHandler() int {
	fmt.Println("bad command called")
	return 1
}
