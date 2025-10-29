/*
A key-value data store, with record tagging.

Use tags to organise, group and search your records.

	| Field | Type     | Comments                  |
	| ----- | -------- | ------------------------- |
	| Key   | string   | Primary key.              |
	| Value | string   | Record value.  Free text. |
	| Tags  | []string | Optional tags.            |

The primary key must be unique.  It cannot be empty.  It cannot contain more than 50 characters.
Other than spaces it cannot contain whitespace.  It cannot start or end with a space.  It must be a
valid UTF8 string.

There are two types of tags:

1. System Tags
System tags are ready-only, and always start with a period.  Examples include:

- .created | The record created timestamp.  RFC3339 format.
- .updated | The record last updated timestamp.  RFC3339 format.
- .deleted | Marks a record as deleted.  Deleted records are hidden unless requested.

2. User Tags
Can contain any combination of lowercase letters, numbers and hyphens.  Tags must be between 1 and
20 characters long.
*/
package tagdb

import "dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"

var (
	dbConnection *db
)

type db struct {
	inMemStore *inMemStore
}

func Connect() *db {
	logger.Info("connecting to db")

	if dbConnection == nil {
		logger.Info("initialising db")
		dbConnection = &db{
			inMemStore: connectToInMemStore(),
		}
	}

	return dbConnection
}

// Closes the connection and releases any held resources.
func (db *db) Close() error {
	panic("not implemented")
}

// List records by tags.
// Tags are optional.  When not provided, all records are returned.
// When provided, only records matching all tags are returned.
func (db *db) List(tags []string) ([]TaggedKV, error) {
	panic("not implemented")
}

// Retrieves a record by its key.
func (db *db) Get(key string) (taggedKv TaggedKV, found bool, err error) {
	panic("not implemented")
}

// Creates or updates a record.
func (db *db) Set(key, value string) error {
	panic("not implemented")
}

// Removes a record from the database.
func (db *db) Delete(key string) error {
	panic("not implemented")
}

// Adds a tag to a record.
func (db *db) Tag(key string, tag string) error {
	panic("not implemented")
}

// Removes a tag from a record.
func (db *db) Untag(key string, tag string) error {
	panic("not implemented")
}
