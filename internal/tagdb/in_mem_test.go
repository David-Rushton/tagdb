package tagdb

import (
	"cmp"
	"slices"
	"testing"
)

func Test_InMemStore_ShouldRoundTrip(t *testing.T) {
	// Arrange
	txId := "test-transaction-id"
	testCases := []struct {
		op       operator
		expected TaggedKV
	}{
		{op: &setOperation{transactionId: txId, key: "key1", value: "value1"}, expected: TaggedKV{Key: "key1", Value: "value1"}},
		{op: &setOperation{transactionId: txId, key: "key2", value: "value2"}, expected: TaggedKV{Key: "key2", Value: "value2"}},
		{op: &setOperation{transactionId: txId, key: "key3", value: "value3"}, expected: TaggedKV{Key: "key3", Value: "value3"}},
	}
	db := newInMemStore()

	// Act
	for _, testCase := range testCases {
		db.apply([]operator{testCase.op})
	}

	// Assert
	for _, testCase := range testCases {
		actual, found := db.get(testCase.expected.Key)
		if !found {
			t.Errorf("Expected to find key `%s`, but it was not found.", testCase.expected.Key)
			continue
		}

		if actual.Key != testCase.expected.Key {
			t.Errorf("Expected key `%s`, but got `%s`.",
				testCase.expected.Key, actual.Key)
		}

		if actual.Value != testCase.expected.Value {
			t.Errorf("Expected value `%s` for key `%s`, but got `%s`.",
				testCase.expected.Value, testCase.expected.Key, actual.Value)
		}
	}
}

func Test_InMemStore_get_ShouldNotReturnUnknownKeys(t *testing.T) {
	// Arrange
	db := newInMemStore()
	ops := []operator{
		&setOperation{transactionId: "tx1", key: "key1", value: "value1"},
	}
	db.apply(ops)

	// Act
	taggedKV, found := db.get("invalid-key")

	// Assert
	if found {
		t.Errorf("Invalid 'invalid-key' found with key `%s` and value `%s`", taggedKV.Key, taggedKV.Value)
	}
}

func Test_InMemStore_get_ShouldNotReturnDeletedKeys(t *testing.T) {
	// Arrange
	db := newInMemStore()
	ops := []operator{
		&setOperation{transactionId: "tx1", key: "to-be-deleted", value: "value1"},
		&deleteOperation{transactionId: "tx1", key: "deleted-key"},
	}
	db.apply(ops)

	// Act
	taggedKV, found := db.get("deleted-key")

	// Assert
	if found {
		t.Errorf("Deleted 'deleted-key' found with key `%s` and value `%s`", taggedKV.Key, taggedKV.Value)
	}
}

func Test_InMemStore_list_ShouldReturnTaggedKeys(t *testing.T) {
	// Arrange
	db := newInMemStore()
	ops := []operator{
		&setOperation{transactionId: "tx1", key: "key-1", value: "value-1"},
		&setOperation{transactionId: "tx1", key: "key-2", value: "value-2"},
		&setOperation{transactionId: "tx1", key: "key-3", value: "value-3"},
		&tagOperation{transactionId: "tx1", key: "key-1", tag: "find-me"},
		&tagOperation{transactionId: "tx2", key: "key-2", tag: "find-me"},
	}
	db.apply(ops)

	// Act
	taggedKVs := db.list([]string{"find-me"})

	// Assert
	expectedCount := 2
	if len(taggedKVs) != expectedCount {
		t.Errorf("Expected %d tagged KVs, but found %d", expectedCount, len(taggedKVs))
	}

	slices.SortFunc(taggedKVs, sortTaggedKVs)

	if taggedKVs[0].Key != "key-1" {
		t.Errorf("Expected first tagged KV to have key 'key-1', but got '%s'", taggedKVs[0].Key)
	}

	if taggedKVs[0].Value != "value-1" {
		t.Errorf("Expected first tagged KV to have value 'value-1', but got '%s'", taggedKVs[0].Value)
	}

	if taggedKVs[1].Key != "key-2" {
		t.Errorf("Expected second tagged KV to have key 'key-2', but got '%s'", taggedKVs[1].Key)
	}

	if taggedKVs[1].Value != "value-2" {
		t.Errorf("Expected second tagged KV to have value 'value-2', but got '%s'", taggedKVs[0].Value)
	}
}

func Test_InMemStore_list_ShouldNotReturnUntaggedKeys(t *testing.T) {
	// Arrange
	db := newInMemStore()
	ops := []operator{
		&setOperation{transactionId: "tx1", key: "key-1", value: "value-1"},
		&setOperation{transactionId: "tx1", key: "key-2", value: "value-2"},
		&setOperation{transactionId: "tx1", key: "key-3", value: "value-3"},
		&tagOperation{transactionId: "tx1", key: "key-1", tag: "find-me"},
		&tagOperation{transactionId: "tx2", key: "key-2", tag: "find-me"},
		&untagOperation{transactionId: "tx2", key: "key-2", tag: "find-me"},
	}
	db.apply(ops)

	// Act
	taggedKVs := db.list([]string{"find-me"})

	// Assert
	expectedCount := 1
	if len(taggedKVs) != expectedCount {
		t.Errorf("Expected %d tagged KVs, but found %d", expectedCount, len(taggedKVs))
	}

	slices.SortFunc(taggedKVs, sortTaggedKVs)

	if taggedKVs[0].Key != "key-1" {
		t.Errorf("Expected first tagged KV to have key 'key-1', but got '%s'", taggedKVs[0].Key)
	}

	if taggedKVs[0].Value != "value-1" {
		t.Errorf("Expected first tagged KV to have value 'value-1', but got '%s'", taggedKVs[0].Value)
	}
}

func sortTaggedKVs(a, b TaggedKV) int {
	return cmp.Compare(a.Key, b.Key)
}

func Test_InMemStore_list_ShouldReturnAllKeysWhenNoTagsProvided(t *testing.T) {
	// Arrange
	db := newInMemStore()
	ops := []operator{
		&setOperation{transactionId: "tx1", key: "key-1", value: "value-1"},
		&setOperation{transactionId: "tx1", key: "key-2", value: "value-2"},
		&setOperation{transactionId: "tx1", key: "key-3", value: "value-3"},
	}
	db.apply(ops)

	// Act
	taggedKVs := db.list([]string{})

	// Assert
	expectedCount := 3
	if len(taggedKVs) != expectedCount {
		t.Errorf("Expected %d tagged KVs, but found %d", expectedCount, len(taggedKVs))
	}
}

func Test_InMemStore_list_ShouldReturnNothingWhenNoneMatch(t *testing.T) {
	// Arrange
	db := newInMemStore()
	ops := []operator{
		&setOperation{transactionId: "tx1", key: "key-1", value: "value-1"},
		&setOperation{transactionId: "tx1", key: "key-2", value: "value-2"},
		&setOperation{transactionId: "tx1", key: "key-3", value: "value-3"},
		&tagOperation{transactionId: "tx1", key: "key-1", tag: "some-tag"},
		&tagOperation{transactionId: "tx2", key: "key-2", tag: "another-tag"},
	}
	db.apply(ops)

	// Act
	taggedKVs := db.list([]string{"non-existent-tag"})

	// Assert
	expectedCount := 0
	if len(taggedKVs) != expectedCount {
		t.Errorf("Expected %d tagged KVs, but found %d", expectedCount, len(taggedKVs))
	}
}

func Test_InMemStore_list_ShouldReturnNothingWhenNoneMatchAll(t *testing.T) {
	// Arrange
	db := newInMemStore()
	ops := []operator{
		&setOperation{transactionId: "tx1", key: "key-1", value: "value-1"},
		&setOperation{transactionId: "tx1", key: "key-2", value: "value-2"},
		&setOperation{transactionId: "tx1", key: "key-3", value: "value-3"},
		&tagOperation{transactionId: "tx1", key: "key-1", tag: "tag-1"},
		&tagOperation{transactionId: "tx2", key: "key-2", tag: "tag-2"},
	}
	db.apply(ops)

	// Act
	taggedKVs := db.list([]string{"tag-1", "tag-2"})

	// Assert
	expectedCount := 0
	if len(taggedKVs) != expectedCount {
		t.Errorf("Expected %d tagged KVs, but found %d", expectedCount, len(taggedKVs))
	}
}

func Test_InMemStore_list_ShouldReturnItemsWithMatchingTags(t *testing.T) {
	// Arrange
	db := newInMemStore()
	ops := []operator{
		&setOperation{transactionId: "tx1", key: "key-1", value: "value-1"},
		&setOperation{transactionId: "tx1", key: "key-2", value: "value-2"},
		&setOperation{transactionId: "tx1", key: "key-3", value: "value-3"},
		&tagOperation{transactionId: "tx1", key: "key-1", tag: "tag-1"},
		&tagOperation{transactionId: "tx1", key: "key-1", tag: "tag-2"},
		&tagOperation{transactionId: "tx1", key: "key-2", tag: "tag-1"},
		&tagOperation{transactionId: "tx1", key: "key-2", tag: "tag-2"},
		&tagOperation{transactionId: "tx2", key: "key-2", tag: "tag-3"},
		&tagOperation{transactionId: "tx2", key: "key-3", tag: "tag-3"},
	}
	db.apply(ops)

	// Act
	taggedKVs := db.list([]string{"tag-1", "tag-2"})

	// Assert
	expectedCount := 2
	if len(taggedKVs) != expectedCount {
		t.Errorf("Expected %d tagged KVs, but found %d", expectedCount, len(taggedKVs))
	}

	slices.SortFunc(taggedKVs, sortTaggedKVs)

	if taggedKVs[0].Key != "key-1" {
		t.Errorf("Expected first tagged KV to have key 'key-1', but got '%s'", taggedKVs[0].Key)
	}

	if taggedKVs[1].Key != "key-2" {
		t.Errorf("Expected second tagged KV to have key 'key-2', but got '%s'", taggedKVs[0].Key)
	}
}
