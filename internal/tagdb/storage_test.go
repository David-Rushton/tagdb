package tagdb

import (
	"cmp"
	"slices"
	"testing"
)

func Test_storage_list_ReturnsItemsWithTag(t *testing.T) {
	// Arrange.
	store, err := openStorage(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to connect to storage: %v", err)
	}
	defer store.close()

	if err := store.set("key-1", "value-1"); err != nil {
		t.Fatalf("set returned error: %s", err)
	}
	if err := store.set("key-2", "value-2"); err != nil {
		t.Fatalf("set returned error: %s", err)
	}
	if err := store.set("key-3", "value-1"); err != nil {
		t.Fatalf("set returned error: %s", err)
	}
	if err := store.tag("key-1", "find"); err != nil {
		t.Fatalf("tag returned error: %s", err)
	}
	if err := store.tag("key-2", "find"); err != nil {
		t.Fatalf("tag returned error: %s", err)
	}

	// Act.
	items, err := store.list([]string{"find"})
	if err != nil {
		t.Fatalf("list returned error: %v", err)
	}

	// Assert.
	var expectedCount = 2
	if len(items) != expectedCount {
		t.Fatalf("list returned %d items but expected %d", len(items), expectedCount)
	}
}

func Test_storage_get_ReturnsItem(t *testing.T) {
	// Arrange.
	store, err := openStorage(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to connect to storage: %v", err)
	}
	defer store.close()

	if err := store.set("key-1", "value-1"); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	// Act.
	taggedKV, found, err := store.get("key-1")

	// Assert.
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	if !found {
		t.Fatalf("Get did not find the item")
	}

	if taggedKV.Key != "key-1" || taggedKV.Value != "value-1" {
		t.Fatalf("Get returned incorrect item: %+v", taggedKV)
	}
}

func Test_storage_delete_RemovesItem(t *testing.T) {
	// Arrange.
	store, err := openStorage(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to connect to storage: %v", err)
	}
	defer store.close()

	if err := store.set("key-1", "value-1"); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}
	if err := store.delete("key-1"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	// Act.
	_, found, err := store.get("key-1")

	// Assert.
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	if found {
		t.Fatalf("Get unexpectedly found the item %+v", found)
	}
}

func Test_storage_list_DoesNotReturnsUntaggedItems(t *testing.T) {
	// Arrange.
	store, err := openStorage(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to connect to storage: %v", err)
	}
	defer store.close()

	if err := store.set("key-1", "value-1"); err != nil {
		t.Fatalf("set returned error: %s", err)
	}
	if err := store.set("key-2", "value-2"); err != nil {
		t.Fatalf("set returned error: %s", err)
	}
	if err := store.set("key-3", "value-1"); err != nil {
		t.Fatalf("set returned error: %s", err)
	}
	if err := store.tag("key-1", "find"); err != nil {
		t.Fatalf("tag returned error: %s", err)
	}
	if err := store.tag("key-2", "find"); err != nil {
		t.Fatalf("tag returned error: %s", err)
	}
	if err := store.untag("key-1", "find"); err != nil {
		t.Fatalf("untag returned error: %s", err)
	}
	if err := store.untag("key-2", "find"); err != nil {
		t.Fatalf("untag returned error: %s", err)
	}

	// Act.
	items, err := store.list([]string{"find"})
	if err != nil {
		t.Fatalf("list returned error: %v", err)
	}

	// Assert.
	var expectedCount = 0
	if len(items) != expectedCount {
		t.Fatalf("list returned %d items but expected %d", len(items), expectedCount)
	}
}

func Test_openStorage_RetainsDataAfterReopen(t *testing.T) {
	// Arrange.
	storeRoot := t.TempDir()
	store, err := openStorage(storeRoot)
	if err != nil {
		t.Fatalf("Failed to connect to storage: %v", err)
	}

	store.set("key 1", "value 1")
	store.set("key 2", "value 2")
	store.set("key 3", "value 3")
	store.tag("key 1", "tag 1")
	store.tag("key 2", "tag 2")
	store.tag("key 3", "tag 3")
	store.untag("key 3", "tag 3")

	store.close()

	store, err = openStorage(storeRoot)
	if err != nil {
		t.Fatalf("Failed to reconnect to storage: %v", err)
	}
	defer store.close()

	// Act.
	items, err := store.list([]string{})
	if err != nil {
		t.Fatalf("list returned error: %v", err)
	}
	slices.SortFunc(items, func(a, b TaggedKV) int {
		return cmp.Compare(a.Key, b.Key)
	})

	// Assert.
	if len(items) != 3 {
		t.Fatalf("Expected 3 items after reopen, but found %d", len(items))
	}

	if items[0].Key != "key 1" || items[0].Value != "value 1" {
		t.Fatalf("Item 1 mismatch after reopen: %+v", items[0])
	}

	if len(items[0].Tags) != 1 || items[0].Tags[0] != "tag 1" {
		t.Fatalf("Item 1 tags mismatch after reopen: %+v", items[0].Tags)
	}

	if items[1].Key != "key 2" || items[1].Value != "value 2" {
		t.Fatalf("Item 2 mismatch after reopen: %+v", items[1])
	}

	if len(items[1].Tags) != 1 || items[1].Tags[0] != "tag 2" {
		t.Fatalf("Item 2 tags mismatch after reopen: %+v", items[0].Tags)
	}

	if items[2].Key != "key 3" || items[2].Value != "value 3" {
		t.Fatalf("Item 3 mismatch after reopen: %+v", items[2])
	}

	if len(items[2].Tags) != 0 {
		t.Fatalf("Item 3 tags mismatch after reopen: %+v", items[0].Tags)
	}
}
