package bimap_test

import (
	"slices"
	"testing"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/bimap"
)

func Test_BiMap_ShouldFindValuesByKey(t *testing.T) {
	bm := &bimap.BiMap[string]{}
	bm.Add("key1", "value1")
	bm.Add("key1", "value2")
	bm.Add("key2", "value3")

	values := bm.GetValues("key1")
	slices.Sort(values)
	expectedValues := []string{"value1", "value2"}

	if slices.Compare(values, expectedValues) != 0 {
		t.Errorf("expected %v, got %v", expectedValues, values)
	}
}

func Test_BiMap_ShouldFindKeysByValue(t *testing.T) {
	bm := &bimap.BiMap[int]{}
	bm.Add(1, 1)
	bm.Add(2, 1)
	bm.Add(3, 2)

	values := bm.GetKeys(1)
	slices.Sort(values)
	expectedValues := []int{1, 2}

	if slices.Compare(values, expectedValues) != 0 {
		t.Errorf("expected %v, got %v", expectedValues, values)
	}
}

func Test_BiMap_ShouldNotReturnDeletedKeysValuePairs(t *testing.T) {
	bm := &bimap.BiMap[int]{}
	bm.Add(1, 1)
	bm.Add(1, 2) // Will be removed.
	bm.Add(1, 3)
	bm.Remove(1, 2)

	values := bm.GetValues(1)
	slices.Sort(values)
	expectedValues := []int{1, 3}

	if slices.Compare(values, expectedValues) != 0 {
		t.Errorf("expected %v, got %v", expectedValues, values)
	}
}
