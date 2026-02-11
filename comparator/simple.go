package comparator

import "github.com/google/go-cmp/cmp"

// SimpleComparator compares status codes, headers, and body using go-cmp.
type SimpleComparator struct{}

func NewSimpleComparator() Comparator {
	return &SimpleComparator{}
}

func (s *SimpleComparator) Compare(newer, current *Response) *Result {
	diff := cmp.Diff(current, newer)

	return &Result{
		Match:      diff == "",
		Newer:      newer,
		Current:    current,
		Difference: diff,
	}
}
