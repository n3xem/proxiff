package proxy

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/n3xem/proxiff/comparator"
)

// Proxy represents the HTTP proxy that forwards requests to newer and current servers
type Proxy struct {
	newerURL   string
	currentURL string
	comparator comparator.Comparator
	client     *http.Client
	logger     *slog.Logger
}

// NewProxy creates a new Proxy instance
func NewProxy(newerURL, currentURL string, comp comparator.Comparator, logger *slog.Logger) *Proxy {
	if logger == nil {
		logger = slog.Default()
	}
	return &Proxy{
		newerURL:   newerURL,
		currentURL: currentURL,
		comparator: comp,
		client:     &http.Client{},
		logger:     logger,
	}
}

// ServeHTTP handles incoming HTTP requests
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the request body once
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		p.logger.Error("failed to read request body",
			slog.String("error", err.Error()),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Forward request to newer server
	newerResp, err := p.forwardRequest(p.newerURL, r, bodyBytes)
	if err != nil {
		p.logger.Error("failed to forward request to newer server",
			slog.String("error", err.Error()),
			slog.String("server", p.newerURL),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		http.Error(w, "Failed to forward to newer server", http.StatusBadGateway)
		return
	}

	// Forward request to current server
	currentResp, err := p.forwardRequest(p.currentURL, r, bodyBytes)
	if err != nil {
		p.logger.Error("failed to forward request to current server",
			slog.String("error", err.Error()),
			slog.String("server", p.currentURL),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		http.Error(w, "Failed to forward to current server", http.StatusBadGateway)
		return
	}

	// Compare responses
	result := p.comparator.Compare(newerResp, currentResp)

	// Log comparison result
	if !result.Match {
		p.logger.Warn("response difference detected",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("newer_status", newerResp.StatusCode),
			slog.Int("current_status", currentResp.StatusCode),
			slog.String("difference", result.Difference),
		)
	} else {
		p.logger.Info("responses match",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", currentResp.StatusCode),
		)
	}

	// Return current response to client
	p.writeResponse(w, currentResp)
}

// forwardRequest forwards the request to the specified URL
func (p *Proxy) forwardRequest(baseURL string, originalReq *http.Request, bodyBytes []byte) (*comparator.Response, error) {
	// Parse target URL
	targetURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	// Build full URL with path and query
	targetURL.Path = originalReq.URL.Path
	targetURL.RawQuery = originalReq.URL.RawQuery

	// Create new request
	req, err := http.NewRequest(originalReq.Method, targetURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	// Copy headers
	for key, values := range originalReq.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Send request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Create comparator.Response
	return &comparator.Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
		Body:       respBody,
	}, nil
}

// writeResponse writes the response to the client
func (p *Proxy) writeResponse(w http.ResponseWriter, resp *comparator.Response) {
	// Copy headers
	for key, values := range resp.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Write status code
	w.WriteHeader(resp.StatusCode)

	// Write body
	w.Write(resp.Body)
}
