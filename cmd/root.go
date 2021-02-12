package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "wilson",
	Short: "wilson - quickly write tests for APIs using JSON configs",
	Long:  `wilson - quickly write tests for APIs using JSON configs`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
