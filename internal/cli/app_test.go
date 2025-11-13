package cli

import "testing"

func Test_app_Run_InvokesExpectedCommand(t *testing.T) {
	// Arrange
	builder := Builder{}
	branch, err := builder.AddBranch("test", "test commands")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = branch.AddCommand("foo", "foo command", func() int {
		return 42
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = branch.AddCommand("bar", "bar command", func() int {
		return -1
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	app := builder.Build()
	app.exit = func(code int) {
		if code != 42 {
			// Assert
			t.Fatalf("expected exit code 42, got %d", code)
		}
	}

	// Act
	app.Run([]string{"my_app", "test", "foo"})
}

func Test_app_Run_Exits1_WhenCommandNotFound(t *testing.T) {
	// Arrange
	builder := Builder{}
	branch, err := builder.AddBranch("test", "test commands")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = branch.AddCommand("foo", "foo command", func() int {
		return 42
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = branch.AddCommand("bar", "bar command", func() int {
		return -1
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	app := builder.Build()
	app.exit = func(code int) {
		if code != -1 {
			// Assert
			t.Fatalf("expected exit code -1, got %d", code)
		}
	}

	// Act
	app.Run([]string{"my_app", "baz"})
}

func Test_app_Run_Exits1_WhenInvalidNumberOfArgsProvided(t *testing.T) {
	// Arrange
	builder := Builder{}
	app := builder.Build()

	actual := 0
	app.exit = func(code int) {
		actual = code
	}

	// Act / Assert
	app.Run([]string{"my_app"})
	if actual != -1 {
		t.Fatalf("expected exit code -1, got %d", actual)
	}

	app.Run([]string{})
	if actual != -1 {
		t.Fatalf("expected exit code -1, got %d", actual)
	}
}
