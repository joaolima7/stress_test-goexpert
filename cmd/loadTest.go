package cmd

import (
	"fmt"
	"time"

	"github.com/joaolima7/stress_test-goexpert/internal/loadtest"
	"github.com/spf13/cobra"
)

var (
	url         string
	requests    int
	concurrency int
)

var loadTestCmd = &cobra.Command{
	Use:   "load",
	Short: "Run load test against a URL",
	Long:  `Execute a load test with specified number of requests and concurrency level.`,
	RunE:  runLoadTest,
}

func init() {
	loadTestCmd.Flags().StringVar(&url, "url", "", "URL to test (required)")
	loadTestCmd.Flags().IntVar(&requests, "requests", 100, "Total number of requests")
	loadTestCmd.Flags().IntVar(&concurrency, "concurrency", 10, "Number of concurrent requests")

	loadTestCmd.MarkFlagRequired("url")
}

func runLoadTest(cmd *cobra.Command, args []string) error {
	if requests <= 0 {
		return fmt.Errorf("requests must be greater than 0")
	}
	if concurrency <= 0 {
		return fmt.Errorf("concurrency must be greater than 0")
	}
	if concurrency > requests {
		concurrency = requests
	}

	fmt.Printf("Starting load test...\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Total Requests: %d\n", requests)
	fmt.Printf("Concurrency: %d\n", concurrency)
	fmt.Println("---")

	config := loadtest.Config{
		URL:         url,
		Requests:    requests,
		Concurrency: concurrency,
		Timeout:     30 * time.Second,
	}

	tester := loadtest.NewloaderTester(config)
	report, err := tester.Run()
	if err != nil {
		return fmt.Errorf("load test failed: %w", err)
	}

	printReport(report)
	return nil
}

func printReport(report *loadtest.Report) {
	fmt.Println("\n=== LOAD TEST REPORT ===")
	fmt.Printf("Total execution time: %v\n", report.TotalTime)
	fmt.Printf("Total requests made: %d\n", report.TotalRequests)
	fmt.Printf("Successful requests (200): %d\n", report.SuccessfulRequests)
	fmt.Printf("Failed requests: %d\n", report.FailedRequests)

	if len(report.StatusCodes) > 0 {
		fmt.Println("\nStatus code distribution:")
		for code, count := range report.StatusCodes {
			fmt.Printf("  %d: %d requests\n", code, count)
		}
	}

	if report.TotalRequests > 0 {
		successRate := float64(report.SuccessfulRequests) / float64(report.TotalRequests) * 100
		fmt.Printf("\nSuccess rate: %.2f%%\n", successRate)
		fmt.Printf("Average requests per second: %.2f\n", float64(report.TotalRequests)/report.TotalTime.Seconds())
	}
}
