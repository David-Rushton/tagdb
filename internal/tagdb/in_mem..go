package tagdb

import (
	"dev.azure.com/trayport/Hackathon/_git/Q/internal/bimap"
	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
)

type inMemStore struct {
	data  map[string]string
	index bimap.BiMap[string]
}

func newInMemStore() *inMemStore {
	logger.Info("initializing in-mem store")

	return &inMemStore{
		data:  map[string]string{},
		index: bimap.BiMap[string]{},
	}
}

// List records by tags.
// Tags are optional.  When not provided, all records are returned.
// When provided, only records matching all tags are returned.
func (db *inMemStore) list(tags []string) []TaggedKV {
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

func (db *inMemStore) apply(op []operator) {
	logger.Infof("applying %d operation(s) to in-mem store", len(op))
	for _, operation := range op {
		switch o := operation.(type) {
		case *setOperation:
			logger.Infof("applying in-mem set operation: key=`%s`, value=`%s`", o.key, o.value)
			db.data[o.key] = o.value

		case *deleteOperation:
			logger.Infof("applying in-mem delete operation: key=`%s`", o.key)
			delete(db.data, o.key)

		case *tagOperation:
			logger.Infof("applying in-mem tag operation: key=`%s`, tag=`%s`", o.key, o.tag)
			db.index.Add(o.key, o.tag)

		case *untagOperation:
			logger.Infof("applying in-mem untag operation: key=`%s`, tag=`%s`", o.key, o.tag)
			db.index.Remove(o.key, o.tag)

		case *commitOperation:
			// No-op.

		default:
			// The in-mem store **must** never diverge from the wal.
			// There is no recovery mechanism.
			logger.Panicf("cannot apply unknown operation type: %+v", operation)
		}
	}
}
