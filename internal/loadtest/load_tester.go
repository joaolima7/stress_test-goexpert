package loadtest

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type LoadTester struct {
	config Config
	client *http.Client
}

type result struct {
	statusCode int
	success    bool
	err        error
}

func NewloaderTester(cfg Config) *LoadTester {
	return &LoadTester{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (lt *LoadTester) Run() (*Report, error) {
	startTime := time.Now()

	results := make(chan result, lt.config.Requests)

	semaphore := make(chan struct{}, lt.config.Concurrency)

	var wg sync.WaitGroup

	for i := 0; i < lt.config.Requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			res := lt.makeRequest()
			results <- res
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	report := &Report{
		StatusCodes: make(map[int]int),
	}

	for res := range results {
		report.TotalRequests++

		if res.success {
			report.SuccessfulRequests++
		} else {
			report.FailedRequests++
		}

		if res.statusCode > 0 {
			report.StatusCodes[res.statusCode]++
		}
	}

	report.TotalTime = time.Since(startTime)

	return report, nil
}

func (lt *LoadTester) makeRequest() result {
	ctx, cancel := context.WithTimeout(context.Background(), lt.config.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", lt.config.URL, nil)
	if err != nil {
		return result{err: err}
	}

	resp, err := lt.client.Do(req)
	if err != nil {
		return result{err: err}
	}
	defer resp.Body.Close()

	success := resp.StatusCode == http.StatusOK

	return result{
		statusCode: resp.StatusCode,
		success:    success,
	}
}
