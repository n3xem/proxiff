package comparator

import (
	"net/http"
	"testing"
)

func TestSimpleComparator_Compare_ExactMatch(t *testing.T) {
	comp := NewSimpleComparator()

	newer := &Response{
		StatusCode: 200,
		Headers:    http.Header{"Content-Type": []string{"application/json"}},
		Body:       []byte(`{"message":"hello"}`),
	}

	current := &Response{
		StatusCode: 200,
		Headers:    http.Header{"Content-Type": []string{"application/json"}},
		Body:       []byte(`{"message":"hello"}`),
	}

	result := comp.Compare(newer, current)

	if !result.Match {
		t.Errorf("Expected match=true, got match=false")
	}

	if result.Difference != "" {
		t.Errorf("Expected no difference, got: %s", result.Difference)
	}
}

func TestSimpleComparator_Compare_StatusCodeDifference(t *testing.T) {
	comp := NewSimpleComparator()

	newer := &Response{
		StatusCode: 200,
		Headers:    http.Header{"Content-Type": []string{"application/json"}},
		Body:       []byte(`{"message":"hello"}`),
	}

	current := &Response{
		StatusCode: 404,
		Headers:    http.Header{"Content-Type": []string{"application/json"}},
		Body:       []byte(`{"message":"hello"}`),
	}

	result := comp.Compare(newer, current)

	if result.Match {
		t.Errorf("Expected match=false, got match=true")
	}

	if result.Difference == "" {
		t.Errorf("Expected difference description, got empty string")
	}
}

func TestSimpleComparator_Compare_BodyDifference(t *testing.T) {
	comp := NewSimpleComparator()

	newer := &Response{
		StatusCode: 200,
		Headers:    http.Header{"Content-Type": []string{"application/json"}},
		Body:       []byte(`{"message":"hello"}`),
	}

	current := &Response{
		StatusCode: 200,
		Headers:    http.Header{"Content-Type": []string{"application/json"}},
		Body:       []byte(`{"message":"world"}`),
	}

	result := comp.Compare(newer, current)

	if result.Match {
		t.Errorf("Expected match=false, got match=true")
	}

	if result.Difference == "" {
		t.Errorf("Expected difference description, got empty string")
	}
}

func TestSimpleComparator_Compare_HeaderDifference(t *testing.T) {
	comp := NewSimpleComparator()

	newer := &Response{
		StatusCode: 200,
		Headers:    http.Header{"Content-Type": []string{"application/json"}},
		Body:       []byte(`{"message":"hello"}`),
	}

	current := &Response{
		StatusCode: 200,
		Headers:    http.Header{"Content-Type": []string{"text/plain"}},
		Body:       []byte(`{"message":"hello"}`),
	}

	result := comp.Compare(newer, current)

	if result.Match {
		t.Errorf("Expected match=false, got match=true")
	}

	if result.Difference == "" {
		t.Errorf("Expected difference description, got empty string")
	}
}
