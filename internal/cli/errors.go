package cli

import "fmt"

// Command names must tbe unique under within a single branch.
type DuplicateCommandError struct {
	CommandName string
}

func (e *DuplicateCommandError) Error() string {
	return fmt.Sprintf("duplicate command `%s`", e.CommandName)
}

type InvalidCommandNameError struct {
	Reason string
}
