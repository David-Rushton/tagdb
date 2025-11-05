package tagdb

import (
	"fmt"
	"math"
	"os"
)

func fileExists(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf("cannot check file existence because %s", err)
	}

	return true, nil
}

func createDirIfNotExists(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, 0755)
		}

		return err
	}

	return nil
}

func createFileIfNotExists(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			file, createErr := os.Create(filePath)
			if createErr != nil {
				return createErr
			}
			return file.Close()
		}
	}

	return nil
}

// Tests if subset is a subset of the superset, or empty.
func isSubsetOrEmpty[T comparable](subset, superset []T) bool {
	if len(subset) == 0 {
		return true
	}

	return isSubset(subset, superset)
}

// Tests if subset is a subset of the superset.
func isSubset[T comparable](subset, superset []T) bool {
	// No need to check, we cannot match them all subset items.
	if len(subset) > len(superset) {
		return false
	}

	for i := range len(subset) {
		var found bool
		for j := range len(superset) {
			if subset[i] == superset[j] {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

type diffResult[T comparable] struct {
	deletes []T
	adds    []T
	inBoth  []T
}

func diff[T comparable](source, target []T) diffResult[T] {
	var result diffResult[T]

	leftMap := toFoundMap(source)
	rightMap := toFoundMap(target)

	for leftKey := range leftMap {
		if _, found := rightMap[leftKey]; found {
			// In both.
			result.inBoth = append(result.inBoth, leftKey)
		} else {
			// Lef only.
			result.deletes = append(result.deletes, leftKey)
		}

		rightMap[leftKey] = false
	}

	for rightKey, found := range rightMap {
		if found {
			// Right only.
			result.adds = append(result.adds, rightKey)
		}
	}

	return result
}

func toFoundMap[T comparable](items []T) map[T]bool {
	foundMap := make(map[T]bool)
	for _, item := range items {
		foundMap[item] = true
	}
	return foundMap
}

// Returns a slice of items that are in both left and right slices.
// Duplicates are preserved.
// Ordering is not guaranteed.
func intersect[T comparable](left, right []T) []T {
	var leftCount = map[T]int{}
	var rightCount = map[T]int{}

	for _, item := range left {
		leftCount[item]++
	}

	for _, item := range right {
		rightCount[item]++
	}

	var result []T
	for item, lCount := range leftCount {
		if rCount, found := rightCount[item]; found {
			minCount := int(math.Min(float64(lCount), float64(rCount)))
			for range minCount {
				result = append(result, item)
			}
		}
	}
	return result
}

func prepend[T any](s []T, v T) []T {
	return append([]T{v}, s...)
}
