package loadtest

import "time"

type Report struct {
	TotalTime          time.Duration
	TotalRequests      int
	SuccessfulRequests int
	FailedRequests     int
	StatusCodes        map[int]int
}
