package tagdb

import (
	"errors"
	"testing"
)

func Test_validateTags_ShouldApproveEmptySlice(t *testing.T) {
	err := validateTags([]string{})
	if err != nil {
		t.Errorf("expected no error for empty slice, but got: %v", err)
	}
}

func Test_validateTags_ShouldApproveValidTags(t *testing.T) {
	err := validateTags([]string{
		"tag-1",
		"tag-2",
		"tag-3",
	})

	if err != nil {
		t.Errorf("expected no error for valid tags: %v", err)
	}
}

func Test_validateTags_ShouldRejectInvalidTags(t *testing.T) {
	err := validateTags([]string{
		"1-tag",                     // cannot start with digit
		"tag_2",                     // underscore not allowed
		"tag@3",                     // special char not allowed
		"TAG-4",                     // uppercase not allowed
		"tag-which-is-way-too-long", // exceeds max length
	})

	if err == nil {
		t.Errorf("expected error for invalid tags, but got none")
	}

	// TODO: Ridiculous error message.
	_ = []error{
		errors.New("tags must match pattern '^[a-z0-9-]{1,20}$'"),
		errors.New("tags must match pattern '^[a-z0-9-]{1,20}$'"),
		errors.New("tags must match pattern '^[a-z0-9-]{1,20}$'"),
		errors.New("tags must match pattern '^[a-z0-9-]{1,20}$'"),
		errors.New("tags must match pattern '^[a-z0-9-]{1,20}$'"),
	}
	// expectedError := errors.Join(expectedErrors...)

	// if err != expectedError {
	// 	t.Errorf("expected error `%s` but got `%s`", err, expectedError)
	// }
}
