package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stress-test",
	Short: "A CLI tool for load testing web services",
	Long:  `A simple but powerful CLI tool to perform load testing on web services with customizable concurrency and request count.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loadTestCmd)
}
