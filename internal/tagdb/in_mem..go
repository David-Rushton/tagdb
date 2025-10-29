package tagdb

import (
	"sync"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/bimap"
	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
)

var (
	inMemStoreConnection *inMemStore
)

type inMemStore struct {
	data  map[string]string
	index bimap.BiMap[string]
	mu    sync.RWMutex
}

func connectToInMemStore() *inMemStore {
	logger.Info("connecting to in-mem store")

	if inMemStoreConnection == nil {
		logger.Info("initializing in-mem store")
		inMemStoreConnection = &inMemStore{}
	}

	return &inMemStore{
		data:  map[string]string{},
		index: bimap.BiMap[string]{},
		mu:    sync.RWMutex{},
	}
}

// List records by tags.
// Tags are optional.  When not provided, all records are returned.
// When provided, only records matching all tags are returned.
func (db *inMemStore) list(tags []string) []TaggedKV {
	db.mu.RLock()
	defer db.mu.RUnlock()

	logger.Infof("in-mem list with tags %v", tags)

	var result []TaggedKV

	if len(tags) == 0 {
		// Return all records.
		for key, value := range db.data {
			recordTags := db.index.GetValues(key)
			result = append(result, TaggedKV{
				Key:   key,
				Value: value,
				Tags:  recordTags,
			})
		}

		return result
	}

	// Find records that match all tags.
	var keysToReturn []string
	for i, tag := range tags {
		taggedKeys := db.index.GetKeys(tag)

		if len(taggedKeys) == 0 {
			// No record matches all tags.  Return empty.
			return []TaggedKV{}
		}

		switch i {
		case 0:
			keysToReturn = taggedKeys

		default:
			keysToReturn = intersect(keysToReturn, taggedKeys)
		}

		if len(keysToReturn) == 0 {
			return []TaggedKV{}
		}
	}

	// Build result.
	for _, key := range keysToReturn {
		value := db.data[key]
		recordTags := db.index.GetValues(key)
		result = append(result, TaggedKV{
			Key:   key,
			Value: value,
			Tags:  recordTags,
		})
	}

	return result
}

// Retrieves a record by its key.
func (db *inMemStore) get(key string) (taggedKv TaggedKV, found bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	logger.Infof("in-mem get with keys %s", key)

	value, found := db.data[key]
	if found {
		tags := db.index.GetValues(key)
		return TaggedKV{
			Key:   key,
			Value: value,
			Tags:  tags,
		}, true
	}
	return TaggedKV{}, false
}

// Creates or updates a record.
func (db *inMemStore) set(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	logger.Infof("in-mem set with key %s with value %s", key, value)

	db.data[key] = value
}

// Removes a record from the database.
func (db *inMemStore) delete(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Delete tags.
	for _, tag := range db.index.GetValues(key) {
		db.index.Remove(key, tag)
	}

	// Delete key-value pair.
	delete(db.data, key)
}

// Adds a tag to a record.
func (db *inMemStore) tag(key string, tag string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, found := db.data[key]; found {
		// Add is idempotent.
		// No need to check if the tag already exists.
		db.index.Add(key, tag)
	}
}

// Removes a tag from a record.
func (db *inMemStore) untag(key string, tag string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, found := db.data[key]; found {
		// Safe to call if tag does not exist.
		db.index.Remove(key, tag)
	}
}
