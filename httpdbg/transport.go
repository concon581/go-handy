// httpdbg/transport.go
package httpdbg

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// DebugTransport is a custom RoundTripper that logs detailed request and response information
type DebugTransport struct {
	// mu prevents concurrent writes to stdout
	mu sync.Mutex
	// Transport is the underlying RoundTripper to use
	Transport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface for detailed logging
func (d *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Use the default transport if none is provided
	transport := d.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	// Clone the request body for logging (as it can only be read once)
	var requestBody []byte
	if req.Body != nil {
		requestBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	}

	// Dump the request details
	d.logRequest(req, requestBody)

	// Perform the actual request
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Clone the response body for logging
	responseBody, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(responseBody))

	// Dump the response details
	d.logResponse(resp, responseBody)

	return resp, nil
}

// logRequest prints detailed information about the outgoing HTTP request
func (d *DebugTransport) logRequest(req *http.Request, body []byte) {
	d.mu.Lock()
	defer d.mu.Unlock()

	fmt.Println("======= HTTP REQUEST =======")
	fmt.Printf("URL: %s %s\n", req.Method, req.URL)

	// Print headers
	for k, v := range req.Header {
		fmt.Printf("%s: %v\n", k, v)
	}

	// Print request body
	if len(body) > 0 {
		fmt.Println("\nBody:")
		fmt.Println(string(body))
	}
	fmt.Println("============================")
}

// logResponse prints detailed information about the incoming HTTP response
func (d *DebugTransport) logResponse(resp *http.Response, body []byte) {
	d.mu.Lock()
	defer d.mu.Unlock()

	fmt.Println("======= HTTP RESPONSE =======")
	fmt.Printf("Status: %s\n", resp.Status)

	// Print headers
	for k, v := range resp.Header {
		fmt.Printf("%s: %v\n", k, v)
	}

	// Print response body
	if len(body) > 0 {
		fmt.Println("\nBody:")
		fmt.Println(string(body))
	}
	fmt.Println("=============================")
}

// NewClient creates an HTTP client with debug logging
func NewClient() *http.Client {
	return &http.Client{
		Transport: &DebugTransport{},
	}
}
