package cli

import (
	"fmt"
	"testing"
)

type starship struct {
	Name           string  `arg:"0:<name>" help:"Name of the starship."`
	Registry       string  `arg:"1:<registry>" help:"Registry number of the starship."`
	Class          string  `option:"--class" help:"Class of the starship."`
	Captain        string  `option:"-c|--captain" help:"Name of the captain."`
	MaximumSpeed   float32 `option:"--max-speed" help:"Maximum speed of the starship in warp factor."`
	CommissionYear int     `option:"--commission-year" help:"Year the starship was commissioned."`
	Decommissioned bool    `option:"--decommissioned" help:"True when the starship has been commissioned."`
}

func (uss *starship) Invoke() int {
	return 0
}

func Test_unmarshalArgs_ReturnsExpectedResult(t *testing.T) {
	// Arrange
	args := []string{
		"Enterprise",
		"NCC-1701-D",
		"--class", "Galaxy",
		"-c", "Jean-Luc Picard",
		"--max-speed", "9.65",
		"--commission-year", "2363",
		"--decommissioned",
	}
	actual := &starship{}
	expected := &starship{
		Name:           "Enterprise",
		Registry:       "NCC-1701-D",
		Class:          "Galaxy",
		Captain:        "Jean-Luc Picard",
		MaximumSpeed:   9.65,
		CommissionYear: 2363,
		Decommissioned: true,
	}

	// Act
	err := unmarshalArgs(args, actual)

	// Assert
	if err != nil {
		t.Errorf("unmarshalArgs() returned an error: %v", err)
	}

	if *expected != *actual {
		t.Errorf("unmarshalArgs() failed \n\tExpected = %+v\n\tActual = %+v", *expected, *actual)
	}
}

func Test_unmarshalArgs_ReturnsExpectedResult_WhenArgsUnordered(t *testing.T) {
	// Arrange
	args := []string{
		"--class", "Galaxy",
		"-c", "Jean-Luc Picard",
		"Enterprise",
		"NCC-1701-D",
		"--decommissioned",
		"--commission-year", "2363",
		"--max-speed", "9.65",
	}
	actual := &starship{}
	expected := &starship{
		Name: "Enterprise",

		Registry:       "NCC-1701-D",
		Class:          "Galaxy",
		Captain:        "Jean-Luc Picard",
		MaximumSpeed:   9.65,
		CommissionYear: 2363,
		Decommissioned: true,
	}

	// Act
	err := unmarshalArgs(args, actual)

	// Assert
	if err != nil {
		t.Errorf("unmarshalArgs() returned an error: %v", err)
	}

	if *expected != *actual {
		t.Errorf("unmarshalArgs() failed \n\tExpected = %+v\n\tActual = %+v", *expected, *actual)
	}
}

func Test_unmarshalArgs_ReturnsDefaultValuesForMissingOptions(t *testing.T) {
	// Arrange
	args := []string{
		"Enterprise",
		"NCC-1701-D",
	}
	actual := &starship{}
	expected := &starship{
		Name: "Enterprise",

		Registry:       "NCC-1701-D",
		Class:          "",
		Captain:        "",
		MaximumSpeed:   0.0,
		CommissionYear: 0,
		Decommissioned: false,
	}

	// Act
	err := unmarshalArgs(args, actual)

	// Assert
	if err != nil {
		t.Errorf("unmarshalArgs() returned an error: %v", err)
	}

	if *expected != *actual {
		t.Errorf("unmarshalArgs() failed \n\tExpected = %+v\n\tActual = %+v", *expected, *actual)
	}
}

func Test_unmarshalArgs_ExpectsStructPointer(t *testing.T) {
	// Arrange
	expected := fmt.Errorf("unmarshal target must be a non-nil pointer")

	// Act & Assert
	err := unmarshalArgs([]string{}, starship{})
	if err == nil || err.Error() != expected.Error() {
		t.Errorf("unmarshalArgs() error \n\tExpected = %v\n\tActual = %v", expected, err)
	}

	err = unmarshalArgs([]string{}, nil)
	if err == nil || err.Error() != expected.Error() {
		t.Errorf("unmarshalArgs() error \n\tExpected = %v\n\tActual = %v", expected, err)
	}

	var someBool bool
	err = unmarshalArgs([]string{}, someBool)
	if err == nil || err.Error() != expected.Error() {
		t.Errorf("unmarshalArgs() error \n\tExpected = %v\n\tActual = %v", expected, err)
	}

	var someInt bool
	expectedErr := fmt.Errorf("unmarshal target must be struct")
	err = unmarshalArgs([]string{}, &someInt)
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("unmarshalArgs() error \n\tExpected = %v\n\tActual = %v", expected, err)
	}
}

func Test_unmarshalArgs_ReturnsErr_WhenAdditionalArgsFound(t *testing.T) {
	// Arrange
	args := []string{
		"Enterprise",
		"NCC-1701-D",
		"--captain", "Jean-Luc Picard",
		"--is-force-user",
		"some-value",
	}
	expected := fmt.Errorf("unexpected args: --is-force-user some-value")

	// Act
	err := unmarshalArgs(args, &starship{})

	// Assert
	if err == nil || err.Error() != expected.Error() {
		t.Errorf("unmarshalArgs() error \n\tExpected = %v\n\tActual = %v", expected, err)
	}
}
