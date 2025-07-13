package loadtest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoadTester_Run(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	config := Config{
		URL:         server.URL,
		Requests:    10,
		Concurrency: 2,
		Timeout:     5 * time.Second,
	}

	tester := NewloaderTester(config)
	report, err := tester.Run()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if report.TotalRequests != 10 {
		t.Errorf("Expected 10 total requests, got %d", report.TotalRequests)
	}

	if report.SuccessfulRequests != 10 {
		t.Errorf("Expected 10 successful requests, got %d", report.SuccessfulRequests)
	}

	if report.StatusCodes[200] != 10 {
		t.Errorf("Expected 10 requests with status 200, got %d", report.StatusCodes[200])
	}
}
