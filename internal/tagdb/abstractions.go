package tagdb

// A key-value pair with tags.
type TaggedKV struct {
	// Primary key.  Must be <= 50 characters.
	Key string

	// Free text value.  Can be empty.
	Value string

	// Organise, group and search your records using optional tags.
	// Tags can only contain lowercase letters, numbers and hyphens.
	// Tags must be between 1 and 20 characters long.
	Tags []string
}
