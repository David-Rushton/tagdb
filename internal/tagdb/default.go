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

import (
	"context"
	"errors"
	"time"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
)

var (
	dbConnection *db
)

type db struct {
	storage   *storage
	config    *dbConfig
	isRunning bool
}

func Start(root string, ctx context.Context, configOptions ...dbConfigurer) {
	if dbConnection != nil {
		logger.Info("tagdb already started")
		return
	}

	logger.Info("starting tagdb")

	// Configure.
	config := &dbConfig{}
	config = WithDefaultConfig()(config)
	for _, configOption := range configOptions {
		config = configOption(config)
	}

	// Ensure storage root exists.
	if err := createDirIfNotExists(root); err != nil {
		logger.Panicf("cannot create database storage because %s", err)
	}

	store, err := openStorage(root)
	if err != nil {
		logger.Panicf("cannot open database storage because %s", err)
	}

	// Create connection.
	dbConnection = &db{
		storage:   store,
		config:    config,
		isRunning: true,
	}

	// Start background maintenance tasks.
	go func(store *storage, config *dbConfig, ctx context.Context) {
		if config.backgroundTaskInterval <= 0 {
			logger.Info("background maintenance tasks disabled")
			return
		}

		ticker := time.NewTicker(config.backgroundTaskInterval)

		select {
		case <-ticker.C:
			if dbConnection == nil || !dbConnection.isRunning {
				logger.Infof("shutting down db")
				Stop()
				return
			}

			logger.Info("running maintenance tasks")
			dbConnection.storage.maybeRoll(config.rollWalAfterBytes)

		case <-ctx.Done():
			logger.Infof("shutting down maintenance tasks")
			Stop()
			return
		}

	}(store, config, ctx)
}

func Stop() {
	if dbConnection == nil {
		logger.Info("tagdb not started")
		return
	}

	if !dbConnection.isRunning {
		logger.Info("tagdb already stopped")
		return
	}

	dbConnection.isRunning = false
	dbConnection.storage.close()
}

func Connect() (*db, error) {
	logger.Info("connecting to db")

	if dbConnection == nil {
		err := logger.Error("database not started")
		return nil, err
	}

	if !dbConnection.isRunning {
		err := logger.Error("database not running")
		return nil, err
	}

	return dbConnection, nil
}

// List records by tags.
// Tags are optional.  When not provided, all records are returned.
// When provided, only records matching all tags are returned.
func (db *db) List(tags []string) ([]TaggedKV, error) {
	logger.Infof("db list records with tags `%+v`", tags)

	// Validation.
	if !db.isRunning {
		err := logger.Error("cannot list because database is not running")
		return []TaggedKV{}, err
	}

	if err := validateTags(tags); err != nil {
		return []TaggedKV{}, err
	}

	return db.storage.list(tags)
}

// Retrieves a record by its key.
func (db *db) Get(key string) (taggedKv TaggedKV, found bool, err error) {
	logger.Infof("db get record with key `%v`", key)

	// Validation.
	if !db.isRunning {
		err := logger.Error("cannot get because database is not running")
		return TaggedKV{}, false, err
	}

	if err := validateKey(key); err != nil {
		return TaggedKV{}, false, err
	}

	return db.storage.get(key)
}

// Creates or updates a record.
func (db *db) Set(key, value string) error {
	logger.Infof("db set record with key `%s` and value `%s`", key, value)

	// Validation.
	var err error

	if !db.isRunning {
		notRunningErr := logger.Error("cannot set because database is not running")
		err = errors.Join(err, notRunningErr)
	}

	if keyErr := validateKey(key); keyErr != nil {
		err = errors.Join(err, keyErr)
	}

	if valueErr := validateValue(key); valueErr != nil {
		err = errors.Join(err, valueErr)
	}

	if err != nil {
		return err
	}

	return db.storage.set(key, value)
}

// Removes a record from the database.
func (db *db) Delete(key string) error {
	logger.Infof("db delete record with key `%s`", key)

	// Validation.
	var err error

	if !db.isRunning {
		notRunningErr := logger.Error("cannot delete because database is not running")
		err = errors.Join(err, notRunningErr)
	}

	if keyErr := validateKey(key); err != nil {
		err = errors.Join(err, keyErr)
	}

	if err != nil {
		return err
	}

	return db.storage.delete(key)
}

// Adds a tag to a record.
func (db *db) Tag(key string, tag string) error {
	logger.Infof("db tag record with key `%s` and tag `%s`", key, tag)

	// Validation.
	var err error

	if !db.isRunning {
		notRunningErr := logger.Error("cannot tag because database is not running")
		err = errors.Join(err, notRunningErr)
	}

	if keyErr := validateKey(key); keyErr != nil {
		err = errors.Join(err, keyErr)
	}

	if tagErr := validateTag(tag); tagErr != nil {
		err = errors.Join(err, tagErr)
	}

	if err != nil {
		return err
	}

	return db.storage.tag(key, tag)
}

// Removes a tag from a record.
func (db *db) Untag(key string, tag string) error {
	logger.Infof("db untag record with key `%s` and tag `%s`", key, tag)

	// Validation.
	var err error

	if !db.isRunning {
		notRunningErr := logger.Error("cannot untag because database is not running")
		err = errors.Join(err, notRunningErr)
	}

	if keyErr := validateKey(key); keyErr != nil {
		err = errors.Join(err, keyErr)
	}

	if tagErr := validateTag(tag); tagErr != nil {
		err = errors.Join(err, tagErr)
	}

	if err != nil {
		return err
	}

	return db.storage.untag(key, tag)
}
