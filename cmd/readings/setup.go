package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"productivity.go/internal/setup"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure Notion credentials",
	Run: func(cmd *cobra.Command, args []string) {
		if err := setup.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}


