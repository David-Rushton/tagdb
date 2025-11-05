package tagdb

import (
	"slices"
	"testing"
)

func Test_isSubset(t *testing.T) {
	var testCases = []struct {
		subset   []string
		superset []string
		expected bool
	}{
		{
			subset:   []string{"a", "b", "c"},
			superset: []string{"a", "b", "c"},
			expected: true,
		},
		{
			subset:   []string{},
			superset: []string{},
			expected: true,
		},
		{
			subset:   []string{},
			superset: []string{"a"},
			expected: true,
		},
		{
			subset:   []string{"a", "b"},
			superset: []string{"a", "b", "c", "d"},
			expected: true,
		},
		{
			subset:   []string{"a", "b", "c"},
			superset: []string{"c", "b", "a"},
			expected: true,
		},
		{
			subset:   []string{"a", "b", "c", "d"},
			superset: []string{"a", "b", "c"},
			expected: false,
		},
		{
			subset:   []string{"a", "b", "c"},
			superset: []string{"a", "b", "x"},
			expected: false,
		},
		{
			subset:   []string{"a"},
			superset: []string{},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		var actual = isSubset(testCase.subset, testCase.superset)
		if actual != testCase.expected {
			t.Errorf("Expected: %v.  Actual: %v.  Subset: %v.  Superset: %v.",
				testCase.expected,
				actual,
				testCase.subset,
				testCase.superset)
		}
	}
}

func Test_diff(t *testing.T) {
	testCases := []struct {
		source   []string
		target   []string
		expected diffResult[string]
	}{
		{
			source: []string{"a", "b", "c"},
			target: []string{"a", "b", "c"},
			expected: diffResult[string]{
				deletes: []string{},
				adds:    []string{},
				inBoth:  []string{"a", "b", "c"},
			},
		},
		{
			source: []string{"c", "b", "a"},
			target: []string{"a", "b", "c"},
			expected: diffResult[string]{
				deletes: []string{},
				adds:    []string{},
				inBoth:  []string{"a", "b", "c"},
			},
		},
		{
			source: []string{"a"},
			target: []string{"b"},
			expected: diffResult[string]{
				deletes: []string{"a"},
				adds:    []string{"b"},
				inBoth:  []string{},
			},
		},
		{
			source: []string{"a", "c", "d", "e", "f"},
			target: []string{"b", "c", "x", "y"},
			expected: diffResult[string]{
				deletes: []string{"a", "d", "e", "f"},
				adds:    []string{"b", "x", "y"},
				inBoth:  []string{"c"},
			},
		},
		{
			source: []string{},
			target: []string{},
			expected: diffResult[string]{
				deletes: []string{},
				adds:    []string{},
				inBoth:  []string{},
			},
		},
	}

	isEqual := func(l, r diffResult[string]) bool {
		slices.Sort(l.deletes)
		slices.Sort(l.adds)
		slices.Sort(l.inBoth)
		slices.Sort(r.deletes)
		slices.Sort(r.adds)
		slices.Sort(r.inBoth)

		return slices.Equal(l.deletes, r.deletes) &&
			slices.Equal(l.adds, r.adds) &&
			slices.Equal(l.inBoth, r.inBoth)
	}

	for _, testCase := range testCases {
		actual := diff(testCase.source, testCase.target)
		if !isEqual(testCase.expected, actual) {
			t.Errorf("Expected: %v.  Actual: %v.  Left: %v.  Right: %v.",
				testCase.expected,
				actual,
				testCase.source,
				testCase.target)
		}
	}
}

func Test_intersect(t *testing.T) {
	testCases := []struct {
		left     []string
		right    []string
		expected []string
	}{
		{
			left:     []string{"a", "b", "c"},
			right:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			left:     []string{"a", "b", "b", "c"},
			right:    []string{"a", "b", "c", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			left:     []string{"a", "b", "b", "b"},
			right:    []string{"a", "b", "b", "c"},
			expected: []string{"a", "b", "b"},
		},
		{
			left:     []string{},
			right:    []string{"a", "b", "c"},
			expected: []string{},
		},
		{
			left:     []string{"c", "b", "a"},
			right:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, testCase := range testCases {
		actual := intersect(testCase.left, testCase.right)
		slices.Sort(actual)
		if !slices.Equal(actual, testCase.expected) {
			t.Errorf("Expected: %v.  Actual: %v.  Left: %v.  Right: %v.",
				testCase.expected,
				actual,
				testCase.left,
				testCase.right)
		}
	}
}

func Test_prepend_PreservesOrdering(t *testing.T) {
	testCases := []struct {
		prepend  int
		to       []int
		expected []int
	}{
		{
			prepend:  1,
			to:       []int{2, 3, 4},
			expected: []int{1, 2, 3, 4},
		},
		{
			prepend:  5,
			to:       []int{},
			expected: []int{5},
		},
		{
			prepend:  0,
			to:       []int{1, 2, 3},
			expected: []int{0, 1, 2, 3},
		},
	}

	for _, testCase := range testCases {
		actual := prepend(testCase.to, testCase.prepend)
		if !slices.Equal(actual, testCase.expected) {
			t.Errorf("Expected: %v.  Actual: %v.  Prepend: %v.  To: %v.",
				testCase.expected,
				actual,
				testCase.prepend,
				testCase.to)
		}
	}
}
