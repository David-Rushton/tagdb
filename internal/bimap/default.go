/*
A bidrectional map (bimap) data structure.

Bimap allows you to look up values by keys, and keys by values.
*/
package bimap

// A bidirectional map.
type BiMap[T comparable] struct {
	keys   map[T]map[T]bool
	values map[T]map[T]bool
}

// Adds a key-value pair to the BiMap.
func (bm *BiMap[T]) Add(key, value T) {
	if bm.keys == nil {
		bm.keys = map[T]map[T]bool{}
	}

	if bm.values == nil {
		bm.values = map[T]map[T]bool{}
	}

	if _, found := bm.keys[key]; !found {
		bm.keys[key] = map[T]bool{}
	}

	if _, found := bm.values[value]; !found {
		bm.values[value] = map[T]bool{}
	}

	bm.keys[key][value] = true
	bm.values[value][key] = true
}

// Removes a key-value pair from the BiMap.
func (bm *BiMap[T]) Remove(key, value T) {
	if values, found := bm.keys[key]; found {
		delete(values, value)
	}

	if len(bm.keys[key]) == 0 {
		delete(bm.keys, key)
	}

	if keys, found := bm.values[value]; found {
		delete(keys, key)
	}

	if len(bm.values[value]) == 0 {
		delete(bm.values, value)
	}
}

// Gets all values associated with a key.
func (bm *BiMap[T]) GetValues(key T) []T {
	var result []T
	if values, found := bm.keys[key]; found {
		for value := range values {
			result = append(result, value)
		}
	}
	return result
}

// Gets all keys associated with a value.
func (bm *BiMap[T]) GetKeys(value T) []T {
	var result []T
	if keys, found := bm.values[value]; found {
		for key := range keys {
			result = append(result, key)
		}
	}
	return result
}
