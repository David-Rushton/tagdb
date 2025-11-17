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

	_, err = branch.AddCommand("good", "a good command", &goodInvoker{})
	if err != nil {
		panic(err)
	}

	_, err = branch.AddCommand("bad", "a bad command", &badInvoker{})
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

type goodInvoker struct {
	Name string `arg:"0:<name>" help:"Name of the good invoker."`
	Age  int    `option:"--age" help:"Age of the good invoker."`
}

func (gi *goodInvoker) Invoke() int {
	fmt.Printf("good: %+v\n", gi)
	return 0
}

type badInvoker struct {
	Name string `arg:"0:<name>" help:"Name of the good invoker."`
	Age  int    `option:"--age" help:"Age of the good invoker."`
}

func (bi *badInvoker) Invoke() int {
	fmt.Printf("bad: %+v\n", bi)
	return 1
}
