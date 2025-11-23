package builtin

import (
	"github.com/google/go-cmp/cmp"
	"github.com/n3xem/proxiff/comparator"
)

// SimpleComparator is a basic implementation that compares status codes, headers, and body
type SimpleComparator struct{}

// NewSimpleComparator creates a new SimpleComparator
func NewSimpleComparator() comparator.Comparator {
	return &SimpleComparator{}
}

// Compare compares two responses and returns the result
func (s *SimpleComparator) Compare(newer, current *comparator.Response) *comparator.Result {
	diff := cmp.Diff(current, newer)

	return &comparator.Result{
		Match:      diff == "",
		Newer:      newer,
		Current:    current,
		Difference: diff,
	}
}
