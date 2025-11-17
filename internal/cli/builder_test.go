package cli

import (
	"errors"
	"testing"
)

func Test_Builder_SetsName(t *testing.T) {
	const expected = "app-name"
	builder := &Builder{}
	builder.Name(expected)

	app := builder.Build()

	actual := app.name
	if actual != expected {
		t.Errorf("unexpected builder name, expected '%s', actual '%s'", expected, actual)
	}
}

func Test_Builder_SetsVersion(t *testing.T) {
	const expected = "0.1.0-test"
	builder := &Builder{}
	builder.Version(expected)

	app := builder.Build()

	actual := app.version
	if actual != expected {
		t.Errorf("unexpected builder version, expected '%s', actual '%s'", expected, actual)
	}
}

func Test_Builder_SetsDescription(t *testing.T) {
	const expected = "some description"
	builder := &Builder{}
	builder.Description(expected)

	app := builder.Build()

	actual := app.description
	if actual != expected {
		t.Errorf("unexpected builder description, expected '%s', actual '%s'", expected, actual)
	}
}

func Test_Builder_SetsMockExit(t *testing.T) {
	builder := &Builder{}

	exitCode := -99
	builder.exit = func(code int) {
		exitCode = code
	}

	app := builder.Build()
	app.exit(42)

	if exitCode != 42 {
		t.Errorf("unexpected exit code, expected 42, actual %d", exitCode)
	}
}
func Test_Builder_AddsBranches(t *testing.T) {
	builder := &Builder{}
	builder.AddBranch("get", "get things")
	builder.AddBranch("list", "list things")

	app := builder.Build()

	if len(app.rootCommand.subcommands) != 2 {
		t.Errorf("unexpected number of subcommands, expected 2, actual %d", len(app.rootCommand.subcommands))
	}

	if _, found := app.rootCommand.subcommands["get"]; !found {
		t.Errorf("expected 'get' branch to be found")
	}

	if _, found := app.rootCommand.subcommands["list"]; !found {
		t.Errorf("expected 'list' branch to be found")
	}
}

func Test_Builder_AddsCommands(t *testing.T) {
	builder := &Builder{}
	builder.AddCommand("get", "get things", &invokeExit0{})
	builder.AddCommand("list", "list things", &invokeExit0{})

	app := builder.Build()

	if len(app.rootCommand.subcommands) != 2 {
		t.Errorf("unexpected number of subcommands, expected 2, actual %d", len(app.rootCommand.subcommands))
	}

	if _, found := app.rootCommand.subcommands["get"]; !found {
		t.Errorf("expected 'get' branch to be found")
	}

	if _, found := app.rootCommand.subcommands["list"]; !found {
		t.Errorf("expected 'list' branch to be found")
	}
}

func Test_Build_ReturnsDuplicateCommandError_WhenAddingDuplicateCommand(t *testing.T) {
	builder := &Builder{}

	_, err1 := builder.AddCommand("get", "get things", &invokeExit0{})
	if err1 != nil {
		t.Errorf("unexpected error when adding first command: %s", err1)
	}

	_, err2 := builder.AddCommand("get", "get things again", &invokeExit0{})
	if err2 == nil {
		t.Errorf("expected error when adding duplicate command, got nil")
	}

	var dupErr *DuplicateCommandError
	if !errors.As(err2, &dupErr) {
		t.Errorf("expected DuplicateCommandError when adding duplicate command, got %T", err2)
	}
}

type invokeExit0 command

func (ie *invokeExit0) Invoke() int {
	return 0
}
