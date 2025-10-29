package tagdb

import (
	"slices"
	"testing"
)

func Test_InMemStore_ShouldRoundTripKeyValuePairs(t *testing.T) {
	// Arrange
	testCases := []struct {
		key   string
		value string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
	}
	db := connectToInMemStore()

	// Act
	for _, testCase := range testCases {
		db.set(testCase.key, testCase.value)
	}

	// Assert
	for _, testCase := range testCases {
		taggedKv, found := db.get(testCase.key)
		if !found {
			t.Errorf("Expected to find key %s, but it was not found.", testCase.key)
			continue
		}

		if taggedKv.Key != testCase.key {
			t.Errorf("Expected key %s, but got %s.",
				testCase.key, taggedKv.Key)
		}

		if taggedKv.Value != testCase.value {
			t.Errorf("Expected value %s for key %s, but got %s.",
				testCase.value, testCase.key, taggedKv.Value)
		}
	}
}

func Test_InMemStore_ShouldNotListKeysWithoutTag(t *testing.T) {
	// Arrange
	db := connectToInMemStore()

	// Act
	taggedKV, found := db.get("key-that-does-not-exist")

	// Assert
	if found {
		t.Errorf("Expected not to find key 'key-that-does-not-exist', but it was found with value %s", taggedKV.Value)
	}
}

func Test_InMemStore_ShouldNotListDeletedKeys(t *testing.T) {
	// Arrange
	db := connectToInMemStore()
	db.set("deleted-key", "value to be deleted")

	// Act
	db.delete("deleted-key")
	taggedKV, found := db.get("deleted-key")

	// Assert
	if found {
		t.Errorf("Expected not to find key 'deleted-key', but it was found with value %s", taggedKV.Value)
	}
}

func Test_InMemStore_ShouldNotListDeletedKeysByTag(t *testing.T) {
	// Arrange
	db := connectToInMemStore()
	db.set("deleted-key", "value")
	db.tag("deleted-key", "tag")

	// Act
	db.delete("deleted-key")
	taggedKVs := db.list([]string{"tag"})

	// Assert
	if len(taggedKVs) != 0 {
		t.Errorf("Expected not to find any keys but found %d", len(taggedKVs))
	}
}

func Test_InMemStore_ShouldNotListDeletedTags(t *testing.T) {
	// Arrange
	db := connectToInMemStore()
	db.set("key", "value")
	db.tag("key", "tag")

	// Act
	db.untag("key", "tag")
	taggedKVs := db.list([]string{"tag"})

	// Assert
	if len(taggedKVs) != 0 {
		t.Errorf("Expected not to find any keys but found %d", len(taggedKVs))
	}
}

func Test_InMemStore_ShouldListTaggedKVs(t *testing.T) {
	// Arrange
	db := connectToInMemStore()
	db.set("foo", "find-me")
	db.tag("foo", "find")
	db.set("bar", "and-me")
	db.tag("bar", "find")
	db.set("baz", "but-not-me")

	// Act
	taggedKVs := db.list([]string{"find"})
	slices.SortFunc(taggedKVs, sortTaggedKVs)

	// Assert
	if len(taggedKVs) != 2 {
		t.Errorf("Expected to find 2 tagged KVs, but found %d", len(taggedKVs))
	}

	if taggedKVs[0].Key != "bar" {
		t.Errorf("Expected first tagged KV to have key 'bar', but got '%s'", taggedKVs[0].Key)
	}

	if taggedKVs[1].Key != "foo" {
		t.Errorf("Expected second tagged KV to have key 'foo', but got '%s'", taggedKVs[1].Key)
	}
}

func Test_InMemStore_ShouldNotFindUntaggedKVs(t *testing.T) {
	// Arrange
	db := connectToInMemStore()
	db.set("foo", "some value")
	db.tag("foo", "some-tag")
	db.set("bar", "some value")
	db.tag("bar", "some-tag")
	db.set("baz", "some value")

	// Act
	taggedKVs := db.list([]string{"non-existent-tag"})

	// Assert
	if len(taggedKVs) != 0 {
		t.Errorf("Expected to find 0 tagged KVs, but found %d", len(taggedKVs))
	}
}

func Test_InMemStore_ShouldNotFindPartialMatches(t *testing.T) {
	// Arrange
	db := connectToInMemStore()
	db.set("foo", "some value") // foo has tags a and b
	db.tag("foo", "a")
	db.tag("foo", "b")
	db.set("bar", "some value") // bar has tags b and c
	db.tag("bar", "b")
	db.tag("bar", "c")
	db.set("bar", "some value") // bar has tags a and c
	db.tag("baz", "a")
	db.tag("baz", "c")

	// Act
	taggedKVs := db.list([]string{"a", "b", "c"})

	// Assert
	if len(taggedKVs) != 0 {
		t.Errorf("Expected to find 0 tagged KVs, but found %d", len(taggedKVs))
	}
}

func Test_InMemStore_ShouldListAllTaggedKVs(t *testing.T) {
	// Arrange
	db := connectToInMemStore()
	db.set("foo", "value-foo")
	db.tag("foo", "tag1")
	db.tag("foo", "tag2")
	db.tag("foo", "tag3")
	db.set("bar", "value-bar")
	db.tag("bar", "tag1")
	db.set("baz", "value-baz")

	// Act
	taggedKVs := db.list([]string{})

	// Assert
	if len(taggedKVs) != 3 {
		t.Errorf("Expected to find 3 tagged KVs, but found %d", len(taggedKVs))
	}
}

func sortTaggedKVs(a, b TaggedKV) int {
	switch {
	case a.Key > b.Key:
		return 1
	case a.Key < b.Key:
		return -1
	default:
		return 0
	}
}
