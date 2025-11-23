package proxy

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/n3xem/proxiff/comparator"
)

// mockComparator is a mock implementation of the Comparator interface for testing
type mockComparator struct {
	called       bool
	returnResult *comparator.Result
}

func (m *mockComparator) Compare(newer, current *comparator.Response) *comparator.Result {
	m.called = true
	if m.returnResult != nil {
		return m.returnResult
	}
	return &comparator.Result{
		Match:      true,
		Newer:      newer,
		Current:    current,
		Difference: "",
	}
}

func TestProxy_ServeHTTP_ForwardsToNewerAndCurrent(t *testing.T) {
	// Setup test servers
	newerCalled := false
	newerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newerCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("newer response"))
	}))
	defer newerServer.Close()

	currentCalled := false
	currentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("current response"))
	}))
	defer currentServer.Close()

	// Create proxy
	mock := &mockComparator{}
	p := NewProxy(newerServer.URL, currentServer.URL, mock, nil)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Execute
	p.ServeHTTP(rec, req)

	// Verify both servers were called
	if !newerCalled {
		t.Errorf("Newer server was not called")
	}
	if !currentCalled {
		t.Errorf("Current server was not called")
	}

	// Verify comparator was called
	if !mock.called {
		t.Errorf("Comparator was not called")
	}
}

func TestProxy_ServeHTTP_ReturnsCurrentResponse(t *testing.T) {
	// Setup test servers
	newerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("newer response"))
	}))
	defer newerServer.Close()

	currentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom", "test-value")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("current response"))
	}))
	defer currentServer.Close()

	// Create proxy
	mock := &mockComparator{}
	p := NewProxy(newerServer.URL, currentServer.URL, mock, nil)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Execute
	p.ServeHTTP(rec, req)

	// Verify response is from current server
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rec.Code)
	}

	body, _ := io.ReadAll(rec.Body)
	if string(body) != "current response" {
		t.Errorf("Expected body 'current response', got '%s'", string(body))
	}

	if rec.Header().Get("X-Custom") != "test-value" {
		t.Errorf("Expected header X-Custom=test-value, got %s", rec.Header().Get("X-Custom"))
	}
}

func TestProxy_ServeHTTP_ForwardsRequestMethod(t *testing.T) {
	receivedMethod := ""
	newerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer newerServer.Close()

	currentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer currentServer.Close()

	mock := &mockComparator{}
	p := NewProxy(newerServer.URL, currentServer.URL, mock, nil)

	req := httptest.NewRequest("POST", "/test", nil)
	rec := httptest.NewRecorder()

	p.ServeHTTP(rec, req)

	if receivedMethod != "POST" {
		t.Errorf("Expected method POST, got %s", receivedMethod)
	}
}

func TestProxy_ServeHTTP_ForwardsRequestBody(t *testing.T) {
	receivedBody := ""
	newerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer newerServer.Close()

	currentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer currentServer.Close()

	mock := &mockComparator{}
	p := NewProxy(newerServer.URL, currentServer.URL, mock, nil)

	requestBody := `{"test":"data"}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(requestBody))
	rec := httptest.NewRecorder()

	p.ServeHTTP(rec, req)

	if receivedBody != requestBody {
		t.Errorf("Expected body %q, got %q", requestBody, receivedBody)
	}
}

func TestProxy_ServeHTTP_ForwardsRequestHeaders(t *testing.T) {
	receivedHeader := ""
	newerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeader = r.Header.Get("X-Test-Header")
		w.WriteHeader(http.StatusOK)
	}))
	defer newerServer.Close()

	currentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer currentServer.Close()

	mock := &mockComparator{}
	p := NewProxy(newerServer.URL, currentServer.URL, mock, nil)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Test-Header", "test-value")
	rec := httptest.NewRecorder()

	p.ServeHTTP(rec, req)

	if receivedHeader != "test-value" {
		t.Errorf("Expected header value 'test-value', got %q", receivedHeader)
	}
}
