package cli

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func validateCommandName(name string) error {
	var errs error

	if strings.ToValidUTF8(name, "") != name {
		errs = errors.Join(errs, fmt.Errorf("command names cannot contain non UTF8 characters"))
	}

	for _, r := range name {
		if !unicode.IsGraphic(r) || unicode.IsSpace(r) {
			errs = errors.Join(errs, fmt.Errorf("command names cannot contain non-printable or whitespace characters"))
			break
		}
	}

	return errs
}

func validateDescription(description string) error {
	if strings.ToValidUTF8(description, "") != description {
		return fmt.Errorf("command descriptions cannot contain non UTF8 characters")
	}

	return nil
}
