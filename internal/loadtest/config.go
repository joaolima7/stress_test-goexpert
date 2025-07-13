package loadtest

import "time"

type Config struct {
	URL         string
	Requests    int
	Concurrency int
	Timeout     time.Duration
}
