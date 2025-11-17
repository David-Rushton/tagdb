package cli

import "testing"

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

func Test_Unmarshal_ReturnsExpectedResult(t *testing.T) {
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
		t.Errorf("Unmarshal() returned an error: %v", err)
	}

	if *expected != *actual {
		t.Errorf("Unmarshal() failed \n\tExpected = %+v\n\tActual = %+v", *expected, *actual)
	}
}
