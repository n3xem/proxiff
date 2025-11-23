package comparator

import "net/http"

// Response represents an HTTP response with its metadata
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// Result represents the comparison result between two responses
type Result struct {
	Match      bool
	Newer      *Response
	Current    *Response
	Difference string // Human-readable description of differences
}

// Comparator is the interface for comparing two HTTP responses
type Comparator interface {
	// Compare compares newer and current responses and returns the result
	Compare(newer, current *Response) *Result
}
